package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	unsecureFlag := flag.Bool("u", false, "Prepend http instead of https when passing a url without a protocol")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Add a new redirect to the redirects file\n",)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  add_redirect [flags] [path] url\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  [path] (optional) The path to use that will be redirect from\n")
		fmt.Fprintf(os.Stderr, "  url (required) The url to redirect to. If no protocol is specified it will be added\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  add_redirect http://example.com\n")
		fmt.Fprintf(os.Stderr, "  add_redirect -u example.com\n")
		fmt.Fprintf(os.Stderr, "  add_redirect /custompath http://example.com\n")
	}
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage: go run add_redirect [-u] [path] url")
		fmt.Println("Run add_redirect --help for more information")
		os.Exit(1)
	}

	var path, urlString string
	userProvidedPath := false

	if len(args) == 2 {
		path = args[0]
		urlString = args[1]
		userProvidedPath = true
	} else {
		urlString = args[0]
	}

	validatedURL, err := validateURL(urlString, *unsecureFlag)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		os.Exit(1)
	}

	redirects, err := readRedirects("redirects")
	if err != nil {
		panic(err)
	}

	if userProvidedPath {
		if _, exists := redirects[path]; exists {
			fmt.Printf("The path %s is already used.\n", path)
			os.Exit(1)
		}
	} else {
		for {
			path = generateRandomPath(5)
			if _, exists := redirects[path]; !exists {
				break
			}
		}
	}

	redirect := fmt.Sprintf("%s %s\n", path, validatedURL)
	appendToFile("redirects", redirect)

	fmt.Printf("Added new redirect: %s -> %s\n", path, validatedURL)
}

func generateRandomPath(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func validateURL(urlString string, unsecure bool) (string, error) {
	if _, err := url.ParseRequestURI(urlString); err != nil {
		scheme := "https://"
		if unsecure {
			scheme = "http://"
		}
		newURL := scheme + urlString
		if _, err := url.ParseRequestURI(newURL); err != nil {
			return "", err
		}
		return newURL, nil
	}
	return urlString, nil
}

func appendToFile(filename, text string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		panic(err)
	}
}

func readRedirects(filename string) (map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	redirects := make(map[string]bool)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 2 {
			redirects[parts[0]] = true
		}
	}
	return redirects, scanner.Err()
}
