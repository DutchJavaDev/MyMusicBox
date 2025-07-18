// @ts-nocheck
import { writable } from "svelte/store";
import Home from "../pages/Home.svelte";
import NotFound from "../pages/NotFound.svelte";
import Playlist from "../pages/Playlist.svelte";
import { getSearchParameters, createSearchParameters } from "../scripts/util";

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
export function initializeRouteService() {
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


