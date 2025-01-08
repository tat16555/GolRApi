package main

import (
	"encoding/json" // ใช้สำหรับการจัดการ JSON เช่น การเข้ารหัสและถอดรหัส
	"fmt"           // ใช้สำหรับพิมพ์ข้อความไปยัง console หรือ HTTP response
	"log"           // ใช้สำหรับการบันทึกข้อความ log
	"net/http"      // ใช้สำหรับการสร้าง HTTP server และจัดการคำขอ
	"time"          // ใช้สำหรับการจัดการเวลา เช่น การเพิ่ม Timestamp
)

// Struct สำหรับจัดการเส้นทาง API
type apiHandler struct{}

// Struct สำหรับถอดรหัส JSON จากคำขอ
type PingRequest struct {
	Message string `json:"message"` // ฟิลด์ message ที่จะรับจาก JSON
}

// Struct สำหรับตอบกลับ JSON
type PingResponse struct {
	Message   string `json:"message"`   // ฟิลด์ message สำหรับตอบกลับ
	Timestamp string `json:"timestamp"` // Timestamp ที่แสดงเวลาปัจจุบัน
}

// ฟังก์ชันสำหรับจัดการเส้นทาง /api/
func (apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// ตรวจสอบว่า Method เป็น POST หรือไม่
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // ส่งสถานะ 405 Method Not Allowed
		jsonResponse(w, map[string]string{
			"error":   "Method not allowed", // ข้อความแสดงข้อผิดพลาด
			"allowed": "POST",               // ระบุ Method ที่รองรับ
		}, http.StatusMethodNotAllowed)
		return // ออกจากฟังก์ชันทันที
	}

	// ถอดรหัส JSON จาก Body ของคำขอ
	var pingReq PingRequest
	if err := json.NewDecoder(req.Body).Decode(&pingReq); err != nil {
		jsonResponse(w, map[string]string{
			"error": "Invalid JSON format", // ข้อความแสดงข้อผิดพลาดกรณี JSON ไม่ถูกต้อง
		}, http.StatusBadRequest)
		return // ออกจากฟังก์ชันทันที
	}

	// ตรวจสอบว่า message ว่างหรือไม่
	if pingReq.Message == "" {
		jsonResponse(w, map[string]string{
			"error": "Message field is required", // ข้อความแจ้งเตือนกรณี message ว่าง
		}, http.StatusBadRequest)
		return
	}

	// บันทึกข้อความที่รับเข้ามาใน log
	log.Printf("Received message: %s", pingReq.Message)

	// สร้างข้อมูลตอบกลับ
	response := PingResponse{
		Message:   "pong",                          // ส่งข้อความ "pong" กลับ
		Timestamp: time.Now().Format(time.RFC3339), // ใช้ Timestamp เวลาปัจจุบัน
	}
	// ส่งข้อมูลตอบกลับในรูปแบบ JSON
	jsonResponse(w, response, http.StatusOK)
}

// ฟังก์ชันสำหรับการตอบกลับ JSON
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")      // กำหนด Content-Type เป็น application/json
	w.WriteHeader(statusCode)                               // กำหนดสถานะ HTTP
	if err := json.NewEncoder(w).Encode(data); err != nil { // เข้ารหัสข้อมูลเป็น JSON และเขียนไปยัง Response
		log.Printf("Failed to write response: %v", err) // บันทึกข้อผิดพลาดใน log หากมีปัญหา
	}
}

func main() {
	mux := http.NewServeMux() // สร้าง HTTP multiplexer สำหรับจัดการเส้นทาง

	// จัดการเส้นทาง /api/
	mux.Handle("/api/", apiHandler{}) // กำหนดให้ใช้ apiHandler สำหรับ /api/

	// จัดการเส้นทาง /
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// ตรวจสอบว่า Path เป็น "/" หรือไม่
		if req.URL.Path != "/" {
			http.NotFound(w, req) // ส่งสถานะ 404 Not Found หาก Path ไม่ใช่ "/"
			return
		}
		fmt.Fprintf(w, "Welcome to the enhanced home page!") // แสดงข้อความต้อนรับในหน้า Home
	})

	// เริ่มต้นเซิร์ฟเวอร์บนพอร์ต 3000
	fmt.Println("Server is running on http://localhost:3000") // แสดงข้อความว่าเซิร์ฟเวอร์กำลังทำงาน
	if err := http.ListenAndServe(":3000", mux); err != nil { // เริ่มฟังคำขอบนพอร์ต 3000
		log.Fatalf("Error starting server: %v", err) // แสดงข้อผิดพลาดหากเซิร์ฟเวอร์เริ่มต้นไม่ได้
	}
}
