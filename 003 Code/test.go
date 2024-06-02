package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
// --- 폴더 이름 설정 필요 (paths.txt)  (x)
// 1-2. 파일 압축 : 하나씩 (파일)압축말고 디렉토리 전체 압축(o)
// 2. 권한 문제(컨테이너 생성하기 위한 이미지 생성이 관리자 권한을 요구로 함)
func buildImage() error {
	fmt.Println("----------------------- Start Build -----------------------")
	err := os.MkdirAll("/Carte/images", 0755)
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

	var tarFileName string
	var imagePaths []string
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME "){
			tarFileName = strings.TrimSpace(strings.TrimPrefix(line, "NAME "))
		} else if strings.HasPrefix(line, "PATH ") {
			imagePaths = append(imagePaths, strings.TrimSpace(strings.TrimPrefix(line, "PATH ")))
		}
	}

	if tarFileName == "" {
		return fmt.Errorf("Tar file name is missing in Cartefile.txt")
	}

	if len(imagePaths) == 0 {
		return fmt.Errorf("No directories specified in Cartefile.txt")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 모든 디렉토리를 하나의 tar.gz 파일로 압축
	imageFilePath := "/Carte/images/" + tarFileName
	err = createImage(imagePaths, imageFilePath)
	// err = createImage(imagePaths, filepath.Join(os.Getenv("HOME"), "Carte", "images", tarFileName))
	if err != nil {
		return err
	}

	// 
	err = extractImage(imageFilePath, "/Carte/images/")
	if err != nil {
		return err
	}

	fmt.Println("---------------------- Image Build Complete ----------------------")
	return nil
}

// 이미지 생성 함수
func createImage(srcDirs []string, dstFile string) error {
	args := []string{"-czvf", dstFile}
	for _, srcDir := range srcDirs {
		absPath, err := filepath.Abs(srcDir)
		if err != nil {
			return err
		}
		args = append(args, "-C", filepath.Dir(absPath), filepath.Base(absPath))
	}
	cmd := exec.Command("tar", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 이미지 추출 함수
func extractImage(tarFile, destDir string) error {
	args := []string{"-xzvf", tarFile, "-C", destDir}
	cmd := exec.Command("tar", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}