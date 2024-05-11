package main

import (
    "fmt"
    // "io"
    "os"
    "os/exec"
    "syscall"
    "path/filepath"
)

// 지우고 생성할때, root권한으로 해서 문제가 아주 많음!!!
// 아무래도 root에 있는 거 잘못 건들인듯 ! (그대로 하면 컴터 날라가니까 조심하기 ,,,,,)

// 이미지 빌드는 문제 없음 !
// 압축해서 하는게 파일 용량부분에서도 좋다고 하지만 너무 시간이 오래걸려서 테스트용으로 바로 사용으로 변경함
// 이미지 빌드 함수
func buildImage() error {
    // 이미지 디렉토리 생성
	fmt.Println("Start Mkdir")
	err := os.MkdirAll("/my_container/rootfs", 0755)
	if err != nil {
		return err
	}

	// 필요한 파일 복사
	fmt.Println("Start CP")
	err = copyFile("/bin/bash", "/my_container/rootfs/")
	if err != nil {
		return err
	}
	err = copyFile("/bin/ls", "/my_container/rootfs/")
	if err != nil {
		return err
	}

	// 이미지 파일 생성
	fmt.Println("Start image create")
	err = createImage("/my_container/rootfs", "/my_container/image.tar")
	if err != nil {
		return err
	}

	fmt.Println("Image build complete.")
	return nil

    // 압축 -- 대기 시간이 너무 오래걸림 ++ 파이프 설정 필요
    // tarCmd := exec.Command("tar", "-C", "/my_container/rootfs", "-c", ".")
    // gzipCmd := exec.Command("gzip")
}

// 파일 복사 함수
func copyFile(src, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}

// 이미지 파일 생성 함수
func createImage(srcDir, dstFile string) error {
	cmd := exec.Command("tar", "-C", srcDir, "-cvf", dstFile, ".")
	return cmd.Run()
}

// 명령어가 잘못 됐을 수도 있음 
// 경로도 다시 따져보고 해보기 , 또 저번에 직접 명령어로 구현했었던 거 활용해서 overlay mount 사용해볼것 !!
// 컨테이너 실행 함수
func runContainer() error {

	fmt.Println("Start running container")
	// cgroup 생성
	err := exec.Command("cgcreate", "-g", "cpu,memory:/my_container").Run()
	if err != nil {
		return err
	}

    fmt.Println("Setting cgroup")
	// cgroup 설정
	err = exec.Command("cgset", "-r", "cpu.cfs_quota_us=100000", "/my_container").Run()
	if err != nil {
		return err
	}

    fmt.Println("Setting cgroup&container")
	// 컨테이너 격리 및 Cgroups 할당
	err = exec.Command("unshare", "--mount", "--pid", "--fork").Run()
	if err != nil {
		return err
	}

    fmt.Println("overlay mount")
	// overlay mount
	err = exec.Command("mount", "--rbind", "/my_container/rootfs", "/").Run()
	if err != nil {
		return err
	}

    fmt.Println("pivotRoot")
	// pivot_root
	err = pivotRoot("/my_container/rootfs")
	if err != nil {
		return err
	}

	// --- 관리자 권한으로 실행하면 여기까지 잘됨 -- 

    fmt.Println("new chroot")
	// 새로운 루트로 chroot
	err = syscall.Chroot("/my_container/rootfs")
	if err != nil {
		return err
	}

	// 컨테이너 내부 프로세스 실행
	cmd := exec.Command("/my_container/rootfs/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// pivot_root 함수
func pivotRoot(newRoot string) error {
	putOld := filepath.Join(newRoot, "/.pivot_root")
	
	// this is dangerous !!!!! 
	// if _, err := os.Stat(putOld); err == nil{
	// 	if err := os.RemoveAll(putOld); err != nil{
	// 		return err
	// 	}
	// }

	if _, err := os.Stat(putOld); err == nil{
		return nil
	}

	if err := syscall.Mount(newRoot, newRoot, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	// 새로운 rootfs로 unmount
	if err := syscall.Mount("", newRoot, "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return err
	}

	// old root를 임시로 저장
	if err := os.Mkdir(putOld, 0700); err != nil {
		return err
	}
	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		return err
	}

	// old root를 삭제
	if err := os.Chdir("/"); err != nil {
		return err
	}
	putOld = "/.pivot_root"
	return syscall.Unmount(putOld, syscall.MNT_DETACH)
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
