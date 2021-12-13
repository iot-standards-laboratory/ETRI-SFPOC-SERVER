window.onload = () => {
  if (!window.WebSocket) {
    alert("No WebSocket!");
    return;
  }
  function connect() {
    let ws = new WebSocket(`ws://${window.location.host}/device/73531420-ae8c-47be-9eff-308d946f3a65`);
    ws.onopen = (e) => {
      console.log("onopen", arguments);
    }

    ws.onclose = () => {
      console.log("onclose", arguments);
    }

    ws.onmessage = function (e) {
      console.log(e.data);
      // addMessage(JSON.parse(e.data));
      console.log(JSON.parse(e.data));
    }
    return ws;
  }

  ws = connect();
}