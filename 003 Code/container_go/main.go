package main

import (
    "fmt"
    // "io"
    "os"
    "os/exec"
    "syscall"
)

// 이미지 빌드 함수
func buildImage() error {
    // 이미지 디렉토리 생성
    fmt.Println("Start Mkdir")
    err := os.MkdirAll("my_container/rootfs", 0755)
    if err != nil {
        return err
    }

    fmt.Println("Start CP")
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
    fmt.Println("Start image create")
    tarCmd := exec.Command("tar", "-C", "my_container/rootfs", "-cvf", "my_container/image.tar", ".")
	err = tarCmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Image build complete.")
	return nil

    // 압축 -- 대기 시간이 너무 오래걸림 ++ 파이프 설정 필요
    // tarCmd := exec.Command("tar", "-C", "my_container/rootfs", "-c", ".")
    // gzipCmd := exec.Command("gzip")
}

// 컨테이너 실행 함수
func runContainer() error {
    // 이미지 파일 로드
    untarCmd := exec.Command("sudo", "tar", "-C", "my_container/rootfs", "-xvf", "my_container/image.tar")
    // untarCmd := exec.Command("tar", "-C", "my_container/rootfs", "-xzf", "my_container/image.tar.gz")
    err := untarCmd.Run()
    if err != nil {
        return err
    }

    fmt.Println("Start running container")
    // cgroup 생성
	cgroupCmd := exec.Command("sudo", "cgcreate", "-g", "cpu,memory:/my_container")
	err = cgroupCmd.Run()
	if err != nil {
		return err
	}

    // cgroup 설정
	cgroupSetCmd := exec.Command("sudo", "cgset", "-r", "cpu.cfs_quota_us=100000", "my_container")
	err = cgroupSetCmd.Run()
	if err != nil {
		return err
	}


    // 사용자 네임스페이스 포함하여 컨테이너 실행
    chrootCmd := exec.Command("sudo", "unshare", "--user", "--yj", "--mount-proc", "--pid", "--fork", "chroot", "my_container/rootfs", "/bin/bash")
    chrootCmd.Stdin = os.Stdin
    chrootCmd.Stdout = os.Stdout
    chrootCmd.Stderr = os.Stderr
    chrootCmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
	}

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
