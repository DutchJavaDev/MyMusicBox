// @ts-nocheck
import { writable } from "svelte/store";
import Home from "../pages/Home.svelte";
import NotFound from "../pages/NotFound.svelte";
import Playlist from "../pages/Playlist.svelte";

const componentsPathMap = new Map([
  ["/404", NotFound],
  ["/Home", Home],
  ["/", Home],
  ["/Playlist", Playlist],
]);

const NotFoundRoutePath = "/404";
const NotFoundPathName = "404";


export let pathName = writable("Home");
export let component = writable(componentsPathMap.get(`/${pathName}`));
export let componentParams = writable(getSearchParameters());

// Initializes the route based on the current URL path and search parameters
// If the path does not exist in the componentsPathMap, it sets the NotFound component
export function initializeRoute() {
  let path = window.location.pathname;
  let parameters = getSearchParameters();

  if (!componentsPathMap.has(path)) {
    component.set(componentsPathMap.get(NotFoundRoutePath));
    componentParams.set({ page: path });
    pathName.set(NotFoundPathName);
    return;
  }

  component.set(componentsPathMap.get(window.location.pathname));
  componentParams.set(parameters);

  if (path === "/") {
    path = "/Home";
  }

  pathName.set(path.split("/")[1]);
}

// Sets the current route and updates the component and parameters accordingly
// If the route does not exist, it sets the NotFound component and parameters
export function navigateTo(newRoute, parameters = null) {
  if (!componentsPathMap.has(newRoute)) {
    component.set(componentsPathMap.get(NotFoundRoutePath));
    componentParams.set({ page: newRoute });
    pathName.set(NotFoundPathName);
  } else {
    component.set(componentsPathMap.get(newRoute));
    if (parameters != null) {
      componentParams.set(parameters);
    }
    pathName.set(newRoute.split("/")[1]);
  }

  let URLSearchParams = createSearchParameters(parameters);

  let url = `${newRoute}${URLSearchParams ? `?${URLSearchParams}` : ""}`;

  if (newRoute !== "/") {
    window.history.pushState({}, "", url);
  } else {
    window.history.replaceState({}, "", url);
  }
}

function getSearchParameters() {
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

function createSearchParameters(params) {
  const searchParams = new URLSearchParams();
  for (const key in params) {
    if (Array.isArray(params[key])) {
      params[key].forEach((value) => searchParams.append(key, value));
    } else {
      searchParams.set(key, params[key]);
    }
  }
  return searchParams.toString();
}
