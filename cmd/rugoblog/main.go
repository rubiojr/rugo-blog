package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed goro.svg
var goroSVG []byte

type Post struct {
	Slug  string
	Title string
	Desc  string
	Date  time.Time
	HTML  template.HTML
	Prev  *PostLink
	Next  *PostLink
}

func (p Post) DateFmt() string {
	return p.Date.Format("January 2, 2006")
}

func (p Post) DateISO() string {
	return p.Date.Format("2006-01-02")
}

type PostLink struct {
	Slug  string
	Title string
}

type IndexData struct {
	Posts []Post
}

type PostData struct {
	Post  Post
	Posts []Post
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var (
	h1Re     = regexp.MustCompile(`(?m)^# (.+)\n*`)
	postFile = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(.+)\.md$`)
)

func run() error {
	postsDir := "posts"
	outDir := "web"

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	md := goldmark.New(
		goldmark.WithExtensions(highlighting.NewTable()),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(postsDir, "[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*.md"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	if len(files) == 0 {
		return fmt.Errorf("no post files (YYYY-MM-DD-*.md) found in %s", postsDir)
	}

	posts := make([]Post, 0, len(files))
	for _, path := range files {
		base := filepath.Base(path)
		m := postFile.FindStringSubmatch(base)
		if m == nil {
			continue
		}

		date, err := time.Parse("2006-01-02", m[1])
		if err != nil {
			return fmt.Errorf("parsing date from %s: %w", base, err)
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		text := string(raw)
		title := extractTitle(text)
		desc := extractDesc(text)
		content := stripFirstParagraph(stripFirstH1(text))

		var buf bytes.Buffer
		if err := md.Convert([]byte(content), &buf); err != nil {
			return fmt.Errorf("converting %s: %w", path, err)
		}

		slug := strings.TrimSuffix(base, ".md")
		posts = append(posts, Post{
			Slug:  slug,
			Title: title,
			Desc:  desc,
			Date:  date,
			HTML:  template.HTML(buf.String()),
		})
	}

	// Sort newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// Wire up prev/next (prev = newer, next = older)
	for i := range posts {
		if i > 0 {
			posts[i].Prev = &PostLink{Slug: posts[i-1].Slug, Title: posts[i-1].Title}
		}
		if i < len(posts)-1 {
			posts[i].Next = &PostLink{Slug: posts[i+1].Slug, Title: posts[i+1].Title}
		}
	}

	// Write embedded mascot logo
	if err := os.WriteFile(filepath.Join(outDir, "goro.svg"), goroSVG, 0o644); err != nil {
		return fmt.Errorf("writing goro.svg: %w", err)
	}

	// Render index
	indexFile, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer indexFile.Close()
	if err := tmpl.ExecuteTemplate(indexFile, "index.html", IndexData{Posts: posts}); err != nil {
		return fmt.Errorf("rendering index: %w", err)
	}

	// Render each post
	for _, p := range posts {
		f, err := os.Create(filepath.Join(outDir, p.Slug+".html"))
		if err != nil {
			return err
		}
		if err := tmpl.ExecuteTemplate(f, "post.html", PostData{Post: p, Posts: posts}); err != nil {
			f.Close()
			return fmt.Errorf("rendering %s: %w", p.Slug, err)
		}
		f.Close()
	}

	fmt.Printf("Built %d posts → %s/\n", len(posts), outDir)
	return nil
}

func extractTitle(s string) string {
	m := h1Re.FindStringSubmatch(s)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	return "Untitled"
}

func extractDesc(s string) string {
	body := stripFirstH1(s)
	body = strings.TrimSpace(body)

	var lines []string
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "```") {
			break
		}
		lines = append(lines, trimmed)
	}

	desc := strings.Join(lines, " ")
	if len(desc) > 200 {
		desc = desc[:200]
		if i := strings.LastIndex(desc, " "); i > 100 {
			desc = desc[:i]
		}
		desc += "…"
	}
	return desc
}

func stripFirstH1(s string) string {
	loc := h1Re.FindStringIndex(s)
	if loc != nil && loc[0] < 100 {
		return s[loc[1]:]
	}
	return s
}

// stripFirstParagraph removes text up to the first blank line or heading.
func stripFirstParagraph(s string) string {
	s = strings.TrimSpace(s)
	lines := strings.SplitN(s, "\n", -1)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "```") {
			return strings.Join(lines[i:], "\n")
		}
	}
	return ""
}
