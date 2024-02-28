package convert

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ObsidianToHugo(obsidianPath, hugoPath string) error {
	// list all markdown files in obsidianPath
	// for each file, convert to hugo format and save in hugoPath
	obsidianFiles, err := listMarkdownFiles(obsidianPath)
	if err != nil {
		return err
	}

	hugoContentPath := path.Join(hugoPath, "content")
	hugoFiles, err := listMarkdownFiles(hugoContentPath)
	if err != nil {
		return err
	}

	if err := removeDeletedObsidianFiles(obsidianFiles, hugoFiles); err != nil {
		return err
	}

	if err := copyObsidianFileToHugo(obsidianFiles, hugoFiles, hugoContentPath); err != nil {
		return err
	}

	if err := normalizeHugoFile(hugoContentPath); err != nil {
		return err
	}

	obsidianConfig, err := ReadObsidianConfig(obsidianPath)
	if err != nil {
		return err
	}

	fmt.Println(obsidianConfig.AttachmentFolderPath)

	// copy attachments
	if err := copyAttachments(path.Join(obsidianPath, obsidianConfig.AttachmentFolderPath), hugoPath); err != nil {
		return err
	}

	return nil
}

func copyAttachments(obsidianAttachmentPath, hugoPath string) error {
	// check if obsidianAttachmentPath exists
	if _, err := os.Stat(obsidianAttachmentPath); os.IsNotExist(err) {
		return nil
	}

	// check if hugoPath exists
	hugoAttachmentPath := path.Join(hugoPath, "static")
	if _, err := os.Stat(hugoAttachmentPath); os.IsNotExist(err) {
		if err := os.MkdirAll(hugoAttachmentPath, os.ModePerm); err != nil {
			return err
		}
	}

	// copy files
	err := filepath.Walk(obsidianAttachmentPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		destPath := strings.ReplaceAll(filePath, obsidianAttachmentPath, hugoAttachmentPath)
		destDir := path.Dir(destPath)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				return err
			}
		}

		if err := copyFile(filePath, destPath); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type ObsidianConfig struct {
	AlwaysUpdateLinks    bool   `json:"alwaysUpdateLinks"`
	NewLinkFormat        string `json:"newLinkFormat"`
	UseMarkdownLinks     bool   `json:"useMarkdownLinks"`
	AttachmentFolderPath string `json:"attachmentFolderPath"`
}

func ReadObsidianConfig(obsidianPath string) (*ObsidianConfig, error) {
	configPath := path.Join(obsidianPath, ".obsidian", "app.json")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &ObsidianConfig{}
	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	return config, nil

}

func normalizeHugoFile(hugoContentPath string) error {
	err := filepath.Walk(hugoContentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		// read file
		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// append front matter

		title := strings.ReplaceAll(info.Name(), ".md", "")
		createdDate := info.ModTime().Format("2006-01-02")
		frontMater := "---\ntitle: " +
			title +
			"\ndate: " + createdDate +
			"\n---\n\n"
		file = append([]byte(frontMater), file...)

		// replace all assets/ with static/
		file = []byte(strings.ReplaceAll(string(file), "assets/", "/"))

		// write file
		if err := os.WriteFile(path, file, os.ModePerm); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func copyObsidianFileToHugo(obsidianFiles, hugoFiles map[string]string, hugoContentPath string) error {
	for relativePath, absPath := range obsidianFiles {
		if _, ok := hugoFiles[relativePath]; ok {
			continue
		}

		destFilePath := path.Join(hugoContentPath, relativePath)
		destDir := path.Dir(destFilePath)
		// check if destDir exists
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				return err
			}
		}

		// copy file
		if err := copyFile(absPath, destFilePath); err != nil {
			return err
		}

	}
	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	if err := destFile.Sync(); err != nil {
		return err
	}

	return nil
}
func removeDeletedObsidianFiles(obsidianFiles, hugoFiles map[string]string) error {
	for relativePath, absPath := range hugoFiles {
		if _, ok := obsidianFiles[relativePath]; !ok {
			if err := os.Remove(absPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func listMarkdownFiles(obsidianPath string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(obsidianPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".md" {
			files[strings.ReplaceAll(path, obsidianPath, "")] = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
