package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	green   = "\033[92m"
	red     = "\033[91m"
	reset   = "\033[0m"
	white   = "\033[97m"
	discord = "\033[38;2;88;101;242m"
	banner  = `
         ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
         ⣿⣿⡟⠁              ⠈⢹⣿⣿
         ⣿⣿⡇                ⢸⣿⣿
         ⣿⣿⡇    ⣴⣾⣵⣶⣶⣾⣿⣦⡄   ⢸⣿⣿
         ⣿⣿⡇  ⢀⣾⣿⣿⢿⣿⣿⣿⣿⣿⣿⡄  ⢸⣿⣿
         ⣿⣿⡇  ⢸⣿⣿⣧⣀⣼⣿⣄⣠⣿⣿⣿  ⢸⣿⣿
         ⣿⣿⡇  ⠘⠻⢷⡯⠛⠛⠛⠛⢫⣿⠟⠛  ⢸⣿⣿
         ⣿⣿⡇                ⢸⣿⣿
         ⣿⣿⣧⡀               ⢸⣿⣿
         ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣆⣸⣿⣿
         ⣿⣿⣿⣿⣿⣿⣿ ` + white + `@1hehaq` + discord + ` ⣿⣿⣿⣿⣿⣿

`
)

type Config struct {
	WebhookURL string `json:"webhook_url"`
}

type DiscordMessage struct {
	Content   string `json:"content"`
	TTS       bool   `json:"tts"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func showBanner() {
	fmt.Print(discord + banner + reset)
}

func init() {
	flag.Usage = func() {
		fmt.Printf("\n%susage%s:\n", red, reset)
		fmt.Printf("  %sdoxo%s [flags] [message]\n\n", discord, reset)
		fmt.Printf("%sFlags%s:\n", red, reset)
		fmt.Printf("  %-12s %s\n", fmt.Sprintf(" %s-%s%s   :", discord, "config", reset), "path to config file (default: ~/.doxo/doxo.json)")
		fmt.Printf("  %-12s %s\n", fmt.Sprintf(" %s-%s%s      :", discord, "txt", reset), "send output as text file")
		fmt.Printf("  %-12s %s\n", fmt.Sprintf(" %s-%s%s    :", discord, "plain", reset), "send as plain text message")
		fmt.Printf("  %-12s %s\n", fmt.Sprintf(" %s-%s%s      :", discord, "tts", reset), "send message with text-to-speech enabled")
		fmt.Printf("  %-12s %s\n", fmt.Sprintf(" %s-%s%s     :", discord, "help", reset), "show help message")
	}
}

func main() {
	configFlag := flag.String("config", "", "path to config file (default: ~/.doxo/doxo.json)")
	txtFlag := flag.Bool("txt", false, "send output as text file")
	plainFlag := flag.Bool("plain", false, "send as plain text message")
	ttsFlag := flag.Bool("tts", false, "send message with text-to-speech enabled")
	helpFlag := flag.Bool("help", false, "show help message")
	flag.Parse()

	if *helpFlag {
		showBanner()
		flag.Usage()
		return
	}

	configPath := *configFlag
	if configPath == "" {
		if path, err := ioutil.ReadFile("/tmp/doxo_path"); err == nil {
			configPath = string(path)
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, white+"["+discord+"doxo"+white+"] "+red+"Error getting home dir: %v\n"+reset, err)
				os.Exit(1)
			}
			configPath = filepath.Join(homeDir, ".doxo", "doxo.json")
		}
	} else {
		ioutil.WriteFile("/tmp/doxo_path", []byte(configPath), 0644)
	}

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, white+"["+discord+"doxo"+white+"] "+red+"Error loading config: %v\n"+reset, err)
		os.Exit(1)
	}

	if config.WebhookURL == "" {
		fmt.Fprintf(os.Stderr, white+"["+discord+"doxo"+white+"] "+red+"no webhook configured!"+reset+"\n")
		os.Exit(1)
	}

	if *configFlag != "" && !*txtFlag && !*plainFlag && len(flag.Args()) == 0 {
		showBanner()
		fmt.Printf(white + "[" + discord + "doxo" + white + "] " + green + "config file loaded successfully!" + reset + "\n")
		return
	}

	var content string

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			content += line + "\n"
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, white+"["+discord+"doxo"+white+"] "+red+"Error reading input: %v\n"+reset, err)
			os.Exit(1)
		}
	} else {
		if *plainFlag || *txtFlag || *ttsFlag {
			content = strings.Join(flag.Args(), " ")
		} else {
			fmt.Fprintf(os.Stderr, white+"["+discord+"doxo"+white+"] "+red+"no input provided. stdin or use doxo flags!\n"+reset)
			os.Exit(1)
		}
	}

	var sendErr error
	if *txtFlag {
		sendErr = sendAsFile(config.WebhookURL, content)
	} else if *plainFlag || *ttsFlag {
		sendErr = sendToDiscord(config.WebhookURL, content, *ttsFlag)
	} else {
		sendErr = sendToDiscord(config.WebhookURL, content, false)
	}

	if sendErr != nil {
		fmt.Fprintf(os.Stderr, white+"["+discord+"doxo"+white+"] "+red+"Error sending to Discord: %v"+reset+"\n", sendErr)
		os.Exit(1)
	}

	fmt.Print(content)
}

func loadConfig(configPath string) (Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			showBanner()
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return Config{}, fmt.Errorf("white["+discord+"doxo"+white+"] "+red+"failed to create config directory: %v"+reset, err)
			}

			template := Config{
				WebhookURL: "",
			}

			templateData, err := json.MarshalIndent(template, "", "    ")
			if err != nil {
				return Config{}, fmt.Errorf("white["+discord+"doxo"+white+"] "+red+"failed to create config template: %v"+reset, err)
			}

			if err := ioutil.WriteFile(configPath, templateData, 0644); err != nil {
				return Config{}, fmt.Errorf("white["+discord+"doxo"+white+"] "+red+"failed to write config template: %v"+reset, err)
			}

			fmt.Printf(white+"["+discord+"doxo"+white+"] "+green+"config template created at: "+reset+"%s\n"+reset, configPath)
			fmt.Printf(white + "[" + discord + "doxo" + white + "] " + green + "please set your webhook in the config file\n" + reset)
			os.Exit(0)
		}
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func sendToDiscord(webhookURL, content string, tts bool) error {
	message := DiscordMessage{
		Content:   content,
		TTS:       tts,
		Username:  "doxo",
		AvatarURL: "https://raw.githubusercontent.com/1hehaq/doxo/refs/heads/main/avatar/doxo.jpeg",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Discord returned status code: %d", resp.StatusCode)
	}

	return nil
}

func sendAsFile(webhookURL, content string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("username", "doxo")
	writer.WriteField("avatar_url", "https://raw.githubusercontent.com/1hehaq/doxo/refs/heads/main/avatar/doxo.jpeg")

	part, err := writer.CreateFormFile("file", "output.txt")
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(content))
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", webhookURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord returned status code: %d", resp.StatusCode)
	}

	return nil
}
