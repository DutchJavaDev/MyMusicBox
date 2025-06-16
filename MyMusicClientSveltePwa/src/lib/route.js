import { writable } from "svelte/store";

export const route = writable("Home");

export function initializeRoute() {
    setRoute(window.location.pathname.split("/")[1] === "" ? "Home" : window.location.pathname.split("/")[1])
}

export function setRoute(newRoute) {
  let currentRoute = window.location.pathname;

  if (currentRoute.startsWith(newRoute) && currentRoute !== "/") {
    // If the current route already matches the new route, do nothing
    return;
  }

  if (newRoute !== "/") {
    window.history.pushState({}, "", newRoute);
  } else {
    window.history.replaceState({}, "", newRoute);
  }
  route.set(newRoute);
}
