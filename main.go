package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/shdkej/synoym-dict/dictionary"
	elastic "github.com/shdkej/synoym-dict/elastic"
	pb "github.com/shdkej/synoym-dict/proto"
	"github.com/tinrab/retry"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
)

var (
	listen     = flag.String("listen", ":8080", "listen address")
	grpcListen = flag.String("grpc listen", ":9090", "listen address")
)

type DictServer struct {
	e *elastic.Elastic
	pb.UnimplementedSynonymDictServer
}

func (s *DictServer) CreateSynonym(ctx context.Context, r *pb.Request) (*pb.Synonym, error) {
	tag := dictionary.Tag{Name: r.GetName()}
	m := structs.Map(tag)
	err := es.Set(m)
	if err != nil {
		log.Println(err)
	}
	log.Println(r)
	return &pb.Synonym{Name: tag.Name}, nil
}

func (s *DictServer) GetAll(ctx context.Context, r *pb.Request) (*httpbody.HttpBody, error) {
	synonym, err := es.GetAll()
	log.Println(synonym)
	if err != nil {
		log.Println(err)
	}

	body, err := json.Marshal(synonym)
	if err != nil {
		log.Println(err)
	}
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        body,
	}, nil
}

func (s *DictServer) GetSynonym(ctx context.Context, r *pb.Request) (*httpbody.HttpBody, error) {
	synonym, err := es.GetSynonym(r.GetName())
	log.Println(synonym)
	if err != nil {
		log.Println(err)
	}

	body, err := json.Marshal(synonym)
	if err != nil {
		log.Println(err)
	}
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        body,
	}, nil
}

func (s *DictServer) Update(ctx context.Context, r *pb.Request) (*pb.Synonym, error) {
	name := r.GetName()
	tag := r.Tags
	err := es.Update(name, tag)
	if err != nil {
		log.Println(err)
	}
	return &pb.Synonym{Name: name, Tags: tag}, nil
}

var es *elastic.Elastic

func main() {
	flag.Parse()

	mapping := `{
  "settings": {
    "analysis": {
      "tokenizer": {
        "nori_user_dict": {
          "type": "nori_tokenizer",
          "decompound_mode": "mixed",
          "user_dictionary": "userdict.txt"
        }
      },
      "analyzer": {
        "korean_analyzer": {
          "filter": [
            "pos_filter_speech", "nori_readingform",
            "lowercase", "synonym", "remove_duplicates"
          ],
          "tokenizer": "nori_user_dict"
        }
      },
      "filter": {
        "synonym" : {
          "type" : "synonym_graph",
          "synonyms_path" : "synonyms.txt"
        },
        "pos_filter_speech": {
          "type": "nori_part_of_speech",
          "stoptags": [
            "E", "J", "SC", "SE", "SF", "SP", "SSC", "SSO", "SY", "VCN", "VCP",
            "VSV", "VX", "XPN", "XSA", "XSN", "XSV"
          ]
        }
      }
    }
  }
}`

	retry.ForeverSleep(2*time.Second, func(_ int) error {
		var err error
		es, err = elastic.CreateElasticsearch()
		if err != nil {
			log.Println("retry", err)
			return err
		}
		es.SetMapping(mapping)
		return nil
	})

	s := &DictServer{}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}

func (s *DictServer) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcServer := grpc.NewServer()
	pb.RegisterSynonymDictServer(grpcServer, &DictServer{})
	log.Printf("start gRPC server on %s port", *grpcListen)
	lis, err := net.Listen("tcp", *grpcListen)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	go grpcServer.Serve(lis)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterSynonymDictHandlerFromEndpoint(ctx, gwmux, *grpcListen, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "proto/openapi/dict.swagger.json")
	})
	mux.Handle("/", gwmux)

	log.Println("running on", *listen)
	return http.ListenAndServe(*listen, allowCORS(mux))
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	w.Header().Set("Access-Control-Allow-Methods", "*")
	return
}
