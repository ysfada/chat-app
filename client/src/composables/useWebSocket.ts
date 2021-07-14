import { ref } from 'vue';

const connections = ref(new Map<string, WebSocket>());

export default function useWebSocket(retry = true) {
  function connect(url: string, protocols?: string | string[]) {
    return new Promise<void>(function (resolve, reject) {
      let ws = connections.value.get(url);
      if (ws && ws.readyState !== ws.CLOSED) {
        return resolve();
      }

      console.info('opening WebSocket connection...');
      connections.value.set(url, new WebSocket(url, protocols));
      ws = connections.value.get(url);

      if (ws) {
        ws.addEventListener('open', (_e) => {
          console.info('WebSocket connection is opened...');
          resolve();
        });

        ws.addEventListener('close', async (e) => {
          console.info('WebSocket connection is closed...');
          if (e.code === 1000) return;

          if (retry) {
            console.info('retrying to open WebSocket connection...');
            await connect(url, protocols).catch((err) => reject(err));
          }
        });

        ws.addEventListener('error', (e) => {
          reject(e);
        });
      }
    });
  }

  return { connect, connections };
}
