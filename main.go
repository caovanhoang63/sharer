package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db *sql.DB
}

type ShareRequest struct {
	HTMLContent string `json:"html_content"`
}

type ShareResponse struct {
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	defer app.db.Close()

	// Routes
	http.HandleFunc("/", app.handleHome)
	http.HandleFunc("/api/share", app.handleShare)
	http.HandleFunc("/shared/", app.handleSharedContent)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NewApp() (*App, error) {
	db, err := sql.Open("sqlite3", "./sharer.db")
	if err != nil {
		return nil, err
	}

	app := &App{db: db}
	if err := app.initDB(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS shared_content (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slug TEXT UNIQUE NOT NULL,
		html_content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_slug ON shared_content(slug);
	`
	_, err := app.db.Exec(query)
	return err
}

func (app *App) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		app.serveHomePage(w, r)
	} else if r.Method == "POST" {
		app.handleFormSubmission(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) serveHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HTML Sharer</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f5f5;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            padding: 40px;
            width: 100%;
            max-width: 800px;
        }
        
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
            font-weight: 600;
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #555;
        }
        
        textarea {
            width: 100%;
            min-height: 200px;
            padding: 12px;
            border: 2px solid #e1e5e9;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            resize: vertical;
            transition: border-color 0.2s;
        }
        
        textarea:focus {
            outline: none;
            border-color: #007bff;
        }
        
        input[type="file"] {
            width: 100%;
            padding: 12px;
            border: 2px dashed #e1e5e9;
            border-radius: 8px;
            background: #f8f9fa;
            cursor: pointer;
            transition: border-color 0.2s;
        }
        
        input[type="file"]:hover {
            border-color: #007bff;
        }
        
        .submit-btn {
            width: 100%;
            padding: 15px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        
        .submit-btn:hover {
            background: #0056b3;
        }
        
        .submit-btn:disabled {
            background: #6c757d;
            cursor: not-allowed;
        }
        
        .result {
            margin-top: 30px;
            padding: 20px;
            border-radius: 8px;
            display: none;
        }
        
        .result.success {
            background: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
        }
        
        .result.error {
            background: #f8d7da;
            border: 1px solid #f5c6cb;
            color: #721c24;
        }
        
        .result-url {
            font-weight: 600;
            word-break: break-all;
        }
        
        .result-url a {
            color: #007bff;
            text-decoration: none;
        }
        
        .result-url a:hover {
            text-decoration: underline;
        }
        
        .loading {
            display: none;
            text-align: center;
            margin-top: 20px;
            color: #666;
        }
        
        @media (max-width: 600px) {
            .container {
                padding: 20px;
            }
            
            textarea {
                min-height: 150px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>HTML Sharer</h1>
        <form id="shareForm">
            <div class="form-group">
                <label for="htmlContent">Paste your HTML code:</label>
                <textarea id="htmlContent" name="htmlContent" placeholder="<!DOCTYPE html>&#10;<html>&#10;<head>&#10;    <title>My Page</title>&#10;</head>&#10;<body>&#10;    <h1>Hello World!</h1>&#10;</body>&#10;</html>"></textarea>
            </div>
            
            <div class="form-group">
                <label for="htmlFile">Or upload an HTML file:</label>
                <input type="file" id="htmlFile" name="htmlFile" accept=".html,.htm">
            </div>
            
            <button type="submit" class="submit-btn" id="submitBtn">Create Link</button>
        </form>
        
        <div class="loading" id="loading">
            Processing your request...
        </div>
        
        <div class="result" id="result">
            <div id="resultContent"></div>
        </div>
    </div>

    <script>
        document.getElementById('shareForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const submitBtn = document.getElementById('submitBtn');
            const loading = document.getElementById('loading');
            const result = document.getElementById('result');
            const resultContent = document.getElementById('resultContent');
            
            // Get form data
            const htmlContent = document.getElementById('htmlContent').value.trim();
            const htmlFile = document.getElementById('htmlFile').files[0];
            
            let contentToSubmit = '';
            
            // Priority: textarea content over file
            if (htmlContent) {
                contentToSubmit = htmlContent;
            } else if (htmlFile) {
                if (!htmlFile.name.toLowerCase().endsWith('.html') && !htmlFile.name.toLowerCase().endsWith('.htm')) {
                    showResult('Please select an HTML file (.html or .htm)', 'error');
                    return;
                }
                
                try {
                    contentToSubmit = await readFileAsText(htmlFile);
                } catch (error) {
                    showResult('Error reading file: ' + error.message, 'error');
                    return;
                }
            } else {
                showResult('Please provide HTML content either by pasting it or uploading a file.', 'error');
                return;
            }
            
            // Show loading state
            submitBtn.disabled = true;
            loading.style.display = 'block';
            result.style.display = 'none';
            
            try {
                const response = await fetch('/api/share', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        html_content: contentToSubmit
                    })
                });
                
                const data = await response.json();
                
                if (response.ok && data.url) {
                    const fullUrl = window.location.origin + data.url;
                    showResult('Your HTML has been shared! <div class="result-url"><a href="' + fullUrl + '" target="_blank">' + fullUrl + '</a></div>', 'success');
                } else {
                    showResult(data.error || 'An error occurred while creating the link.', 'error');
                }
            } catch (error) {
                showResult('Network error: ' + error.message, 'error');
            } finally {
                submitBtn.disabled = false;
                loading.style.display = 'none';
            }
        });
        
        function readFileAsText(file) {
            return new Promise((resolve, reject) => {
                const reader = new FileReader();
                reader.onload = e => resolve(e.target.result);
                reader.onerror = e => reject(new Error('Failed to read file'));
                reader.readAsText(file);
            });
        }
        
        function showResult(message, type) {
            const result = document.getElementById('result');
            const resultContent = document.getElementById('resultContent');
            
            resultContent.innerHTML = message;
            result.className = 'result ' + type;
            result.style.display = 'block';
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func (app *App) handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var htmlContent string

	// Priority: textarea content over file
	textareaContent := r.FormValue("htmlContent")
	if strings.TrimSpace(textareaContent) != "" {
		htmlContent = textareaContent
	} else {
		// Try to read from uploaded file
		file, header, err := r.FormFile("htmlFile")
		if err == nil {
			defer file.Close()

			// Validate file extension
			ext := strings.ToLower(filepath.Ext(header.Filename))
			if ext != ".html" && ext != ".htm" {
				http.Error(w, "Please upload an HTML file", http.StatusBadRequest)
				return
			}

			content, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}
			htmlContent = string(content)
		}
	}

	if strings.TrimSpace(htmlContent) == "" {
		http.Error(w, "No HTML content provided", http.StatusBadRequest)
		return
	}

	// Generate slug and save
	slug := app.generateSlug()
	err = app.saveContent(slug, htmlContent)
	if err != nil {
		http.Error(w, "Error saving content", http.StatusInternalServerError)
		return
	}

	// Redirect to success page or return JSON
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ShareResponse{URL: "/shared/" + slug})
	} else {
		http.Redirect(w, r, "/?success="+slug, http.StatusSeeOther)
	}
}

func (app *App) handleShare(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ShareResponse{Error: "Invalid JSON"})
		return
	}

	if strings.TrimSpace(req.HTMLContent) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ShareResponse{Error: "No HTML content provided"})
		return
	}

	slug := app.generateSlug()
	if err := app.saveContent(slug, req.HTMLContent); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ShareResponse{Error: "Error saving content"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShareResponse{URL: "/shared/" + slug})
}

func (app *App) handleSharedContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/shared/")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	content, err := app.getContent(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			app.serve404(w, r)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

func (app *App) serve404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Page Not Found - HTML Sharer</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f5f5;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            padding: 40px;
            text-align: center;
            max-width: 500px;
        }
        h1 {
            color: #dc3545;
            margin-bottom: 20px;
        }
        p {
            color: #666;
            margin-bottom: 30px;
            line-height: 1.6;
        }
        a {
            color: #007bff;
            text-decoration: none;
            font-weight: 600;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>404 - Page Not Found</h1>
        <p>The shared HTML page you're looking for doesn't exist or may have been removed.</p>
        <a href="/">← Back to HTML Sharer</a>
    </div>
</body>
</html>`
	
	w.Write([]byte(html))
}

func (app *App) generateSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	rand.Seed(time.Now().UnixNano())
	slug := make([]byte, length)
	for i := range slug {
		slug[i] = charset[rand.Intn(len(charset))]
	}
	return string(slug)
}

func (app *App) saveContent(slug, content string) error {
	query := "INSERT INTO shared_content (slug, html_content) VALUES (?, ?)"
	_, err := app.db.Exec(query, slug, content)
	return err
}

func (app *App) getContent(slug string) (string, error) {
	query := "SELECT html_content FROM shared_content WHERE slug = ?"
	var content string
	err := app.db.QueryRow(query, slug).Scan(&content)
	return content, err
}