package main
 
import (
	"encoding/json"
	"log"
	"net/http"
 
	"borrowercopilot/rules"
)
 
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web")))
	mux.HandleFunc("/api/calculate", handleCalculate)
 
	addr := ":8080"
	log.Printf("Borrower Copilot running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
 
func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var in rules.Input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	out := rules.Calculate(in)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}