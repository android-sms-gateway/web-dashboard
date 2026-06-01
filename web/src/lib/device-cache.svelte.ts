import { listDevices } from "./api";
import type { Device } from "./types";

export const dc = $state({
  devices: [] as Device[],
  loading: false,
  error: null as string | null,
});

let loaded = false;
let inflight: Promise<void> | null = null;

export async function loadDevices(force = false) {
  if (loaded && !force) return;
  if (inflight) return inflight;

  dc.loading = true;
  dc.error = null;

  inflight = (async () => {
    try {
      dc.devices = await listDevices();
      loaded = true;
    } catch (e: any) {
      dc.error = e.message ?? "Failed to load devices";
    } finally {
      dc.loading = false;
      inflight = null;
    }
  })();
  return inflight;
}

export function refreshDevices() {
  loaded = false;
  return loadDevices(true);
}
