package cmd

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build command that handles image compression",
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildImage()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

// continaer_3 코드 합치기
// 수정할부분
// 1. 파일 입력(scanf) 부분으로 하지 말고 file내에서 읽어오도록 하기 -> 원하는 경로 복사, 디렉토리 이름 설정 ==> 이미지 압축 생성
// --- file 경로 완료, 해당 폴더 상위에서 Carte build 하면 원하는 경로로 압축 완료(o)
// --- 폴더 이름 설정 필요 (paths.txt)  (o)
// 1-2. 파일 압축 : 하나씩 (파일)압축말고 디렉토리 전체 압축(o)
// 2. 권한 문제(컨테이너 생성하기 위한 이미지 생성이 관리자 권한을 요구로 함) (x)
func buildImage() error {
	fmt.Println("----------------------- Start Build -----------------------")
	baseDir := "/Carte/images"
	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		return err
	}

	// 파일에서 경로를 읽어오기
	file, err := os.Open("Cartefile.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var name string
	var srcPath string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME ") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "NAME "))
		} else if strings.HasPrefix(line, "PATH ") {
			srcPath = strings.TrimSpace(strings.TrimPrefix(line, "PATH "))
		}
	}

	if name == "" {
		return fmt.Errorf("Name is missing in Cartefile.txt")
	}

	if srcPath == "" {
		return fmt.Errorf("Path is missing in Cartefile.txt")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 타겟 디렉토리 경로 설정
	targetDir := filepath.Join(baseDir, name)
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	// 디렉토리 내용을 타겟 디렉토리에 복사
	err = copyDirectoryContents(srcPath, targetDir)
	if err != nil {
		return err
	}

	// 타겟 디렉토리를 tar.gz 파일로 압축
	tarFileName := name + ".tar.gz"
	tarFilePath := filepath.Join(baseDir, tarFileName)
	err = createImage(targetDir, tarFilePath)
	if err != nil {
		return err
	}

	fmt.Println("---------------------- Image Build Complete ----------------------")
	return nil
}

// 디렉토리 내용을 복사하는 함수
func copyDirectoryContents(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

// 파일 복사 함수
func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	info, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(dstFile, info.Mode())
}

// 이미지 생성 함수
func createImage(srcDir, dstFile string) error {
	args := []string{"-czvf", dstFile, "-C", srcDir, "."}
	cmd := exec.Command("tar", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

