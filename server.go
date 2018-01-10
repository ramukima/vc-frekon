package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
    	"net"
	"net/http"
	"io/ioutil"
	"time"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/matryer/way"
)

// Server is the app server.
type Server struct {
	facebox *facebox.Client
	router  *way.Router
}

// NewServer makes a new Server.
func NewServer(facebox *facebox.Client) *Server {
	srv := &Server{
		facebox: facebox,
		router:  way.NewRouter(),
	}
	srv.router.HandleFunc(http.MethodPost, "/faceid", srv.handlewebFaceID)
	srv.router.HandleFunc(http.MethodGet, "/status", srv.handleSystemStatus)
	srv.router.HandleFunc(http.MethodGet, "/arm", srv.handleSystemArm)
	srv.router.HandleFunc(http.MethodGet, "/disarm", srv.handleSystemDisarm)
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handlewebFaceID(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	if err != nil {
        	log.Printf("Form file parameter not found: %v\n", err)
		http.Error(w, "Form file parameter not found", http.StatusInternalServerError)
		return
	}

	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
        	log.Printf("Failed to read uploaded file:%v\n", err)
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}

	faces, err := s.facebox.Check(bytes.NewReader(b))
	if err != nil {
		log.Printf("[ERROR] Error on facebox %v\n", err)
		http.Error(w, "something went wrong verifying the faces", http.StatusInternalServerError)
		return
	}
	var response struct {
		FaceLen int    `json:"faces_len"`
		Matched bool   `json:"matched"`
		Id 		string `json:"id"`
		Name    	string `json:"name"`
		Confidence	float64 `json:"confidence"`
	}
	response.FaceLen = len(faces)
	if len(faces) == 1 {
		response.Matched = faces[0].Matched
		response.Name = faces[0].Name
		response.Id = faces[0].ID
		response.Confidence = faces[0].Confidence
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	b := motionControl("status")
	var response struct {
		Status string `json:"status"`
	}

	response.Status = "System Status: " + b

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSystemArm(w http.ResponseWriter, r *http.Request) {
	b := motionControl("start")
	var response struct {
		Status string `json:"status"`
		Description string `json:"description"`
	}

	response.Status = "System Armed"
	response.Description = "System Armed: " + b

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleSystemDisarm(w http.ResponseWriter, r *http.Request) {
	b := motionControl("pause")
	var response struct {
		Status string `json:"status"`
		Description string `json:"description"`
	}

	response.Status = "System Disarmed"
	response.Description = "System Disarmed: " + b

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func motionControl(action string) string {
	homeURL := os.Getenv("MOTION_URL")
	if homeURL == "" {
		wanIp := GetOutboundIP()
		if wanIp != nil {
			homeURL = fmt.Sprintf("http://%s", wanIp) + ":8082"
		}
	}

	if homeURL == "" {
		homeURL = "http://127.0.0.1:8082"
	}

        timeout := time.Duration(2 * time.Second)
        c := &http.Client{ Timeout: timeout }
        resp, err := c.Get(homeURL + "/0/detection/" + action)
        if err != nil {
                fmt.Printf("http.Get() error: %v\n", err)
                return fmt.Sprintf("%v", err)
        }
        defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return string(body)
		}
		return "Nil Response"
	}
	return "Error from server: " + string(resp.StatusCode)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP
}

