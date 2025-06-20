// @ts-nocheck
import { writable } from "svelte/store";
import Home from "../pages/Home.svelte";
import NotFound from "../pages/NotFound.svelte";
import Playlist from "../pages/Playlist.svelte";

let routes = new Map([
  ["/404", NotFound],
  ["/Home", Home],["/", Home],
  ["/Playlist", Playlist]
]);

// TODO fix route value when path is /
// Its empty string now, but should be "Home" or similar
export let route = writable("Home");
export let component = writable(routes.get("/Home"));
export let componentParams = writable(getSearchParams());

export function initializeRoute() {
  let path = window.location.pathname;
  let parameters = getSearchParams();

  if (!routes.has(path)) {
    component.set(routes.get("/404"));
    componentParams.set({ page: path });
    route.set("404");
    return;
  }

  component.set(routes.get(window.location.pathname));
  componentParams.set(parameters);
  route.set(path.split("/")[1]);
}

export function setRoute(newRoute, parameters) {
  if (!routes.has(newRoute)) {
    component.set(routes.get("/404"));
    componentParams.set({ page: newRoute });
    route.set("404");
  } else {
    component.set(routes.get(newRoute));
    if(parameters != null)
    {
      componentParams.set(parameters);
    }
    route.set(newRoute.split("/")[1]);
  }

  let URLSearchParams = createSearchParams(parameters);

  let url = `${newRoute}${URLSearchParams ? `?${URLSearchParams}` : ""}`;

  if (newRoute !== "/") {
    window.history.pushState({}, "", url);
  } else {
    window.history.replaceState({}, "", url);
  }
}

function getSearchParams() {
  const searchParams = new URLSearchParams(window.location.search);
  const result = {};
  for (const [key, value] of searchParams.entries()) {
    if (result[key]) {
      // If key already exists, convert to array or push into existing array
      result[key] = Array.isArray(result[key]) ? [...result[key], parseValue(value)] : [result[key], parseValue(value)];
    } else {
      result[key] = parseValue(value);
    }
  }
  return result;
}

function parseValue(value) {
  if (value === "true") return true;
  if (value === "false") return false;
  if (!isNaN(value) && value.trim() !== "") return Number(value);
  return value;
}

function createSearchParams(params) {
  const searchParams = new URLSearchParams();
  for (const key in params) {
    if (Array.isArray(params[key])) {
      params[key].forEach(value => searchParams.append(key, value));
    } else {
      searchParams.set(key, params[key]);
    }
  }
  return searchParams.toString();
}