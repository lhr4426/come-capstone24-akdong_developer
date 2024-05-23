package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "syscall"
)

// 이미지 빌드 함수
func buildImage() error {
    // 이미지 디렉토리 생성
    imageDir := "/my_container/rootfs"
    fmt.Println("Creating image directory:", imageDir)
    if err := os.MkdirAll(imageDir, 0755); err != nil {
        return err
    }

    // 필요한 파일 복사
    files := []string{"/bin/bash", "/bin/ls"}
    for _, file := range files {
        fmt.Println("Copying file:", file)
        if err := copyFile(file, imageDir); err != nil {
            return err
        }
    }

    // 이미지 파일 생성
    imageFile := "/my_container/image.tar"
    fmt.Println("Creating image file:", imageFile)
    if err := createImage(imageDir, imageFile); err != nil {
        return err
    }

    fmt.Println("Image build complete.")
    return nil
}

// 파일 복사 함수
func copyFile(src, dstDir string) error {
    cmd := exec.Command("cp", src, dstDir)
    return cmd.Run()
}

// 이미지 파일 생성 함수
func createImage(srcDir, dstFile string) error {
    cmd := exec.Command("tar", "-C", srcDir, "-cvf", dstFile, ".")
    return cmd.Run()
}

// 컨테이너 실행 함수
func runContainer() error {
    fmt.Println("Start running container")

    // cgroup 생성
    fmt.Println("Creating cgroup")
    if err := exec.Command("cgcreate", "-g", "cpu,memory:/my_container").Run(); err != nil {
        return err
    }

    // cgroup 설정
    fmt.Println("Setting cgroup")
    if err := exec.Command("cgset", "-r", "cpu.cfs_quota_us=100000", "/my_container").Run(); err != nil {
        return err
    }

    // memory
    fmt.Println("Setting cgroup")
    if err := exec.Command("cgset", "-r", "memory.limit_in_bytes=209715200", "/my_container").Run(); err != nil {
        return err
    }
    
    // memory
    fmt.Println("Setting cgroup")
    if err := exec.Command("cgset", "-r", "memory.swappiness=0", "/my_container").Run(); err != nil {
        return err
    }

    // 컨테이너 격리 및 Cgroups 할당
    fmt.Println("Setting container isolation")
    if err := syscall.Unshare(syscall.CLONE_NEWNS | syscall.CLONE_NEWPID); err != nil {
        return err
    }

    // overlay mount
    fmt.Println("Overlay mount")
    if err := mountOverlay("/my_container/rootfs", "/"); err != nil {
        return err
    }

    // pivotRoot
    fmt.Println("Pivot root")
    if err := pivotRoot("/tmp/overlay-work/merge"); err != nil {
        return err
    }

    // 새로운 루트로 chroot
    fmt.Println("Chroot")
    if err := syscall.Chroot("."); err != nil {
        return err
    }

    // 컨테이너 내부 프로세스 실행
    fmt.Println("Executing container process")
    cmd := exec.Command("/bin/bash")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return err
    }

    // /proc 파일 시스템 마운트
    if err := mountProc(); err != nil {
        return err
    }

    return nil
}

// proc 파일 시스템 마운트 함수
func mountProc() error {
    fmt.Println("Mounting proc filesystem...")
    if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
        return err
    }
    fmt.Println("Proc filesystem mounted successfully.")
    return nil
}

func mountOverlay(rootfs, mountpoint string) error {
    // 오버레이 마운트를 위한 작업 디렉토리 생성
    workDir := "/tmp/overlay-work"
    fmt.Println("Creating overlay work directory:", workDir)
    if err := os.MkdirAll(workDir, 0755); err != nil {
        return err
    }
    containerDir := "/tmp/overlay-work/container"
    workWorkDir := "/tmp/overlay-work/work"
    mergeDir := "/tmp/overlay-work/merge"
    os.MkdirAll(containerDir, 0755)
    os.MkdirAll(workWorkDir, 0755)
    os.MkdirAll(mergeDir, 0755)

    // 마운트 명령어 실행
    cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", rootfs, containerDir, workWorkDir), mergeDir)
    fmt.Println("Mounting overlay filesystem...")

    // 명령어 실행 및 오류 처리
    if err := cmd.Run(); err != nil {
        // 오류 출력
        fmt.Println("Error mounting overlay filesystem:", err)
        return err
    }

    fmt.Println("Overlay filesystem mounted successfully.")
    return nil
}

func pivotRoot(newRoot string) error {
    putOld := filepath.Join(newRoot, "put_old")

    // 새로운 put_old 디렉터리 생성
    if err := os.MkdirAll(putOld, 0755); err != nil {
        return err
    }

    // 현재 작업 디렉터리를 새로운 루트로 변경
    if err := syscall.Chdir(newRoot); err != nil {
        return err
    }

    cmd := exec.Command("pivot_root", ".", "put_old")
    if err := cmd.Run(); err != nil {
        // 오류 출력
        fmt.Println("@@@@@@@@@@@@@", err)
        return err
    }

    // Pivot 이후 작업 디렉터리를 루트로 변경
    if err := syscall.Chdir("/"); err != nil {
        return err
    }

    // 이전 루트를 unmount
    oldRoot := filepath.Join("/", ".pivot_root")
    if err := syscall.Unmount(oldRoot, syscall.MNT_DETACH); err != nil {
        return err
    }

    // put_old 디렉터리 삭제
    if err := os.RemoveAll(oldRoot); err != nil {
        return err
    }

    return nil
}



func main() {
    // 이미지 빌드
    if err := buildImage(); err != nil {
        fmt.Println("Error building image:", err)
        return
    }

    // 컨테이너 실행
    if err := runContainer(); err != nil {
        fmt.Println("Error running container:", err)
        return
    }
}
