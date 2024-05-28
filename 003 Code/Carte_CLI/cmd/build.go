package cmd

import (
	"bufio"
	"fmt"
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
// 1-2. 파일 압축 : 하나씩 (파일)압축말고 디렉토리 전체 압축
// 2. 권한 문제(컨테이너 생성하기 위한 이미지 생성이 관리자 권한을 요구로 함)
func buildImage() error {
	fmt.Println("----------------------- Start Build -----------------------")

	file, err := os.Open("Cartefile.txt")
	if err != nil {
		return fmt.Errorf("failed to open Cartefile.txt: %w", err)
	}
	defer file.Close()

	tempDir := "/tmp/carte_build"
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var workdir string
	var cmdCommand []string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "COPY":
			if len(parts) != 3 {
				return fmt.Errorf("invalid COPY command in Cartefile.txt: %s", line)
			}
			src := parts[1]
			dst := filepath.Join(tempDir, parts[2])

			err := os.MkdirAll(filepath.Dir(dst), 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(dst), err)
			}

			err = copyDir(src, dst)
			if err != nil {
				return fmt.Errorf("failed to copy from %s to %s: %w", src, dst, err)
			}

		case "WORKDIR":
			if len(parts) != 2 {
				return fmt.Errorf("invalid WORKDIR command in Cartefile.txt: %s", line)
			}
			workdir = filepath.Join(tempDir, parts[1])
			err := os.MkdirAll(workdir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create workdir %s: %w", workdir, err)
			}

		case "CMD":
			if len(parts) < 2 {
				return fmt.Errorf("invalid CMD command in Cartefile.txt: %s", line)
			}
			cmdCommand = parts[1:]
			fmt.Println("CMD command found, will be executed after build.")

		default:
			return fmt.Errorf("unknown command in Cartefile.txt: %s", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading Cartefile.txt: %w", err)
	}

	fmt.Println("----------------------- Start image create -----------------------")
	err = createImage(tempDir, "carte_image.tar")
	if err != nil {
		return fmt.Errorf("failed to create image: %w", err)
	}

	fmt.Println("---------------------- Image Build Complete ----------------------")

	if len(cmdCommand) > 0 {
		fmt.Println("----------------------- Executing CMD command -----------------------")
		cmd := exec.Command(cmdCommand[0], cmdCommand[1:]...)
		cmd.Dir = workdir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run CMD command '%s': %w", strings.Join(cmdCommand, " "), err)
		}
	}

	return nil
}

// 디렉토리 복사 함수
func copyDir(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	return cmd.Run()
}

// 이미지 파일 생성 함수
func createImage(srcDir, dstFile string) error {
	cmd := exec.Command("tar", "-C", srcDir, "-cvf", dstFile, ".")
	return cmd.Run()
}

// 컨테이너 생성 명령어
func runContainer() error {
	fmt.Println("----------------------- Start Create Container -----------------------")
	return nil
}