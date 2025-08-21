package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// 对比两个文件夹中所有子文件的hash值
// 输出只存在于目录a中或目录b中的文件，hash值相同的文件，hash值不同的文件

type FileInfo struct {
	Path string
	Hash string
}

func main() {

	cacheDir, _ := os.UserCacheDir()
	appName := "wwtool"
	folderA := filepath.Join(cacheDir, appName, "KrPcSdk_Mainland_bilibili")
	folderB := filepath.Join(cacheDir, appName, "KrPcSdk_Mainland_official")

	// 收集两个文件夹的文件信息
	filesA, err := collectFiles(folderA)
	if err != nil {
		log.Fatalf("Error collecting files from %s: %v", folderA, err)
	}

	filesB, err := collectFiles(folderB)
	if err != nil {
		log.Fatalf("Error collecting files from %s: %v", folderB, err)
	}

	// 比较文件
	onlyInA, onlyInB, differentHashes, sameHashes := compareFiles(filesA, filesB, folderA, folderB)

	// 写入结果文件
	if err := writeResults(onlyInA, onlyInB, differentHashes, sameHashes); err != nil {
		log.Fatalf("Error writing results: %v", err)
	}

	fmt.Println("Comparison completed. Results saved to:")
	fmt.Println("- only_in_a.txt")
	fmt.Println("- only_in_b.txt")
	fmt.Println("- different_hashes.txt")
	fmt.Println("- same_hashes.txt")
}

// collectFiles 收集文件夹中所有文件的路径和hash值
func collectFiles(root string) (map[string]FileInfo, error) {
	files := make(map[string]FileInfo)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		hash, err := fileHash(path)
		if err != nil {
			return fmt.Errorf("error hashing file %s: %v", path, err)
		}

		// 存储相对路径
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %v", path, err)
		}

		// 统一使用正斜杠，确保跨平台一致性
		relPath = filepath.ToSlash(relPath)

		files[relPath] = FileInfo{
			Path: relPath,
			Hash: hash,
		}

		return nil
	})

	return files, err
}

// fileHash 计算文件的SHA256哈希值
func fileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// compareFiles 比较两个文件夹的文件
func compareFiles(filesA, filesB map[string]FileInfo, pathA, pathB string) (onlyInA, onlyInB, differentHashes, sameHashes []string) {
	// 找出只在A中存在的文件
	for path := range filesA {
		if _, exists := filesB[path]; !exists {
			onlyInA = append(onlyInA, path)
		}
	}

	// 找出只在B中存在的文件
	for path := range filesB {
		if _, exists := filesA[path]; !exists {
			onlyInB = append(onlyInB, path)
		}
	}

	// 找出两个文件夹中都存在但hash不同的文件
	for path, fileA := range filesA {
		if fileB, exists := filesB[path]; exists {
			if fileA.Hash != fileB.Hash {
				differentHashes = append(differentHashes, path)
			} else {
				sameHashes = append(sameHashes, path)
			}
		}
	}

	return onlyInA, onlyInB, differentHashes, sameHashes
}

// writeResults 将结果写入文件
func writeResults(onlyInA, onlyInB, differentHashes, sameHashes []string) error {
	// 写入只在A中存在的文件
	if err := writeToFile("only_in_a.txt", onlyInA); err != nil {
		return fmt.Errorf("error writing only_in_a.txt: %v", err)
	}

	// 写入只在B中存在的文件
	if err := writeToFile("only_in_b.txt", onlyInB); err != nil {
		return fmt.Errorf("error writing only_in_b.txt: %v", err)
	}

	// 写入hash不同的文件
	if err := writeToFile("different_hashes.txt", differentHashes); err != nil {
		return fmt.Errorf("error writing different_hashes.txt: %v", err)
	}

	// 写入hash不同的文件
	if err := writeToFile("same_hashes.txt", sameHashes); err != nil {
		return fmt.Errorf("error writing same_hashes.txt: %v", err)
	}

	return nil
}

// writeToFile 将字符串切片写入文件
func writeToFile(filename string, lines []string) error {
	file, err := os.Create(filepath.Join("tmp", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}

	return nil
}
