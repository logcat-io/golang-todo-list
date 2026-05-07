const todoForm = document.querySelector("#todoForm");
const todoTitle = document.querySelector("#todoTitle");
const todoList = document.querySelector("#todoList");
const emptyState = document.querySelector("#emptyState");
const summary = document.querySelector("#summary");
const message = document.querySelector("#message");
const serverStatus = document.querySelector("#serverStatus");
const refreshButton = document.querySelector("#refreshButton");

let todos = [];

function getField(item, name) {
  return item[name] ?? item[name[0].toUpperCase() + name.slice(1)];
}

function setMessage(text, isError = false) {
  message.textContent = text;
  message.classList.toggle("error", isError);
}

async function request(path, options = {}) {
  const response = await fetch(path, {
    headers: { "Content-Type": "application/json", ...(options.headers ?? {}) },
    ...options,
  });

  if (!response.ok) {
    let text = "요청을 처리하지 못했습니다.";
    try {
      const body = await response.json();
      text = body.error ?? body.Error ?? text;
    } catch {
      text = response.statusText || text;
    }
    throw new Error(text);
  }

  if (response.status === 204) {
    return null;
  }
  return response.json();
}

function renderTodos() {
  todoList.replaceChildren();

  const remaining = todos.filter((item) => !getField(item, "done")).length;
  summary.textContent = `전체 ${todos.length}개 · 남은 일 ${remaining}개`;
  emptyState.hidden = todos.length > 0;

  for (const item of todos) {
    const id = getField(item, "ID");
    const title = getField(item, "Title");
    const done = getField(item, "Done");

    console.log(id);
    console.log(title);
    console.log(done);
    console.log(item);

    const row = document.createElement("li");
    row.className = `todo-item${done ? " done" : ""}`;

    const badge = document.createElement("span");
    badge.className = "todo-id";
    badge.textContent = id;

    const titleNode = document.createElement("span");
    titleNode.className = "todo-title";
    titleNode.textContent = title;

    const completeButton = document.createElement("button");
    completeButton.className = "small";
    completeButton.type = "button";
    completeButton.textContent = done ? "완료됨" : "완료";
    completeButton.disabled = done;
    completeButton.addEventListener("click", () => completeTodo(id));

    const deleteButton = document.createElement("button");
    deleteButton.className = "small danger";
    deleteButton.type = "button";
    deleteButton.textContent = "삭제";
    deleteButton.addEventListener("click", () => deleteTodo(id));

    row.append(badge, titleNode, completeButton, deleteButton);
    todoList.append(row);
  }
}

async function loadTodos() {
  try {
    todos = await request("/todos");
    renderTodos();
    setMessage("");
  } catch (error) {
    setMessage(error.message, true);
  }
}

async function checkHealth() {
  try {
    await request("/health");
    serverStatus.textContent = "서버 연결됨";
    serverStatus.classList.add("ok");
  } catch {
    serverStatus.textContent = "서버 오류";
    serverStatus.classList.remove("ok");
  }
}

async function createTodo(title) {
  await request("/todos", {
    method: "POST",
    body: JSON.stringify({ title }),
  });
  todoTitle.value = "";
  await loadTodos();
}

async function completeTodo(id) {
  try {
    await request(`/todos/${id}/complete`, { method: "PATCH" });
    await loadTodos();
  } catch (error) {
    setMessage(error.message, true);
  }
}

async function deleteTodo(id) {
  try {
    await request(`/todos/${id}`, { method: "DELETE" });
    await loadTodos();
  } catch (error) {
    setMessage(error.message, true);
  }
}

todoForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const title = todoTitle.value.trim();
  if (!title) {
    setMessage("할 일을 입력하세요.", true);
    todoTitle.focus();
    return;
  }

  try {
    await createTodo(title);
  } catch (error) {
    setMessage(error.message, true);
  }
});

refreshButton.addEventListener("click", loadTodos);

checkHealth();
loadTodos();
