## 이미지 빌드를 위한 CLI 생성하기
### Carte를 실행하기 위한 방법
- Carte github다운을 받으면 코드 자체 받을 수 있도록 설정
- CLI로 Carte build와 같이 이미지로 생성하고자 하는 경로에서 명령어 실행
- build는 endpoint로 관련 코드가 실행할 수 있도록 도와줌(Carte는 서버 실행 도구, endpoint가 명령어)
- 경로내에 CarteFile필수로 있어야함
- 이미지 경로는 자동으로 cp하여 rootfs에 저장하고 CarteFile에는 go version과 같은 환경설정 및 디렉토리 이름 설정 가능하도록 설계