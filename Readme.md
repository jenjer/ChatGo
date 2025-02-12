[toc]

### 프로그램 설명
go 언어를 사용해 만드는 기본적인 채팅 프로그램으로
사용하기 위해서는 server, client 각각을 컴파일 해서 원하는 위치에서 사용하면 된다.

다중 채팅을 지원하는 프로그램이다. 여러명이서 id 를 기반으로 통신할 수 있으며
id 설정은 구현중


단, 서버와 클라이언트는 각각의 포트를 열어줘야되며, 방화벽에서 작업을 해주어야된다.

현재 코드는 리눅스 기반으로 동작하도록 구현되어 있으며
윈도우 및 맥에서의 동작을 보장하지 않는다.

## Client 프로그램 동작
실행 시키게 되면 ini파일을 읽어 서버의 ip 주소를 받아온다.
단 ini 보다 우선적으로 argv 를 읽어옴으로
./client 192.168.119.131:8000 과 같은 방식으로 주소와 포트를 설정할 수 있고

동작하지 않을 시
ini 파일을 직접 열어서 수정해도 된다.


## Server 프로그램 동작
그냥 실행시키면 된다. 포트 개방은 잊지말고 진행해야된다.

