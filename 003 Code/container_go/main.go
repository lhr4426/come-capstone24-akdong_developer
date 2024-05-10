package main

import (
    "fmt"
    "os"
    "os/exec"
)

// 이미지 빌드 함수
func buildImage() error {
    // 이미지 디렉토리 생성
    err := os.MkdirAll("my_container/rootfs", 0755)
    if err != nil {
        return err
    }

    // 필요한 파일 복사
    cmd := exec.Command("cp", "/bin/bash", "my_container/rootfs/")
    err = cmd.Run()
    if err != nil {
        return err
    }

    cmd = exec.Command("cp", "/bin/ls", "my_container/rootfs/")
    err = cmd.Run()
    if err != nil {
        return err
    }

    // 이미지 파일 생성
    tarCmd := exec.Command("tar", "-C", "my_container/rootfs", "-c", ".")
    gzipCmd := exec.Command("gzip")
    gzipCmd.Stdin, _ = tarCmd.StdoutPipe()
    gzipCmd.Stdout = os.Stdout

    err = tarCmd.Start()
    if err != nil {
        return err
    }
    err = gzipCmd.Start()
    if err != nil {
        return err
    }
    err = tarCmd.Wait()
    if err != nil {
        return err
    }
    err = gzipCmd.Wait()
    if err != nil {
        return err
    }

    fmt.Println("Image build complete.")
    return nil
}

// 컨테이너 실행 함수
func runContainer() error {
    // 이미지 파일 로드
    untarCmd := exec.Command("tar", "-C", "my_container/rootfs", "-xzf", "my_container/image.tar.gz")
    err := untarCmd.Run()
    if err != nil {
        return err
    }

    // 컨테이너 실행
    chrootCmd := exec.Command("chroot", "my_container/rootfs", "/bin/bash")
    chrootCmd.Stdin = os.Stdin
    chrootCmd.Stdout = os.Stdout
    chrootCmd.Stderr = os.Stderr

    err = chrootCmd.Run()
    if err != nil {
        return err
    }

    fmt.Println("Container run complete.")
    return nil
}

func main() {
    // 이미지 빌드
    err := buildImage()
    if err != nil {
        fmt.Println("Error building image:", err)
        return
    }

    // 컨테이너 실행
    err = runContainer()
    if err != nil {
        fmt.Println("Error running container:", err)
        return
    }
}
