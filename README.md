# todo-cli

Go 입문용 Todo 프로젝트입니다. 하나의 Todo 도메인 코드를 CLI 프로그램과 HTTP API 서버에서 함께 사용합니다. API 서버는 간단한 블랙 앤 화이트 프론트 화면도 같이 제공합니다.

## 구조

```text
todo-cli/
  cmd/
    api/          HTTP API 서버와 프론트 정적 파일 서빙
    todo/         CLI 예제 프로그램
  internal/
    todo/         Todo 모델과 메모리 저장소
  web/            브라우저용 프론트엔드
```

## 요구 사항

- Go 1.25.1 이상

## CLI 실행

```bash
go run ./cmd/todo
```

샘플 Todo 목록을 출력합니다.

## API 서버와 프론트 실행

```bash
go run ./cmd/api
```

기본 포트는 `8080`입니다.

브라우저에서 다음 주소를 엽니다.

```text
http://localhost:8080
```

다른 포트로 실행하려면 `PORT` 환경변수를 사용합니다.

```bash
PORT=18080 go run ./cmd/api
```

이 경우 브라우저 주소는 다음과 같습니다.

```text
http://localhost:18080
```

## API

| Method | Path | 설명 |
|---|---|---|
| `GET` | `/health` | 서버 상태 확인 |
| `GET` | `/todos` | Todo 목록 조회 |
| `POST` | `/todos` | Todo 추가 |
| `PATCH` | `/todos/{id}/complete` | Todo 완료 처리 |
| `DELETE` | `/todos/{id}` | Todo 삭제 |

Todo 추가 요청 예시:

```bash
curl -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"learn go"}'
```

## 테스트

```bash
go test ./...
```

로컬 환경에서 Go 빌드 캐시 권한 문제가 나면 프로젝트 내부 캐시를 지정해서 실행할 수 있습니다.

```bash
GOCACHE="$PWD/.gocache" go test ./...
```

## 메모

- 현재 저장소는 메모리 기반이라 서버를 재시작하면 Todo 데이터가 초기화됩니다.
- 프론트는 별도 빌드 도구 없이 `web/` 디렉터리의 HTML, CSS, JavaScript 파일로 동작합니다.
