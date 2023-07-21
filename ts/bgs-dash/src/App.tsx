import "./App.css";
import {
  NavLink,
  RouterProvider,
  createBrowserRouter,
  useNavigate,
} from "react-router-dom";
import Dash from "./components/Dash/Dash";
import { Disclosure } from "@headlessui/react";
import { Bars3Icon, XMarkIcon } from "@heroicons/react/24/outline";
import Login from "./components/Login/Login";
import { useEffect } from "react";
import Logout from "./components/Logout/Logout";
import Domains from "./components/Domains/Domains";
import Repos from "./components/Repos/Repos";
import Consumers from "./components/Consumers/Consumers";

function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(" ");
}

// Redirect to /login if not authenticated
function RequireAuth({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();

  useEffect(() => {
    if (!localStorage.getItem("admin_route_token")) {
      navigate("/login");
    }
  }, []);

  return children;
}

interface Route {
  path: string;
  name: string;
  element: React.ReactNode;
  requrieAuth?: boolean;
  hideIfAuth?: boolean;
}

const routes: Route[] = [
  {
    path: "/",
    name: "PDS List",
    element: (
      <RequireAuth>
        <Nav />
        <main>
          <div className="mx-auto max-w-7xl px-2 py-6 sm:px-6 lg:px-8">
            <Dash />
          </div>
        </main>
      </RequireAuth>
    ),
    requrieAuth: true,
  },
  {
    path: "/consumers",
    name: "Consumers",
    element: (
      <RequireAuth>
        <Nav />
        <main>
          <div className="mx-auto max-w-7xl px-2 py-6 sm:px-6 lg:px-8">
            <Consumers />
          </div>
        </main>
      </RequireAuth>
    ),
    requrieAuth: true,
  },
  {
    path: "/domain_bans",
    name: "Domain Bans",
    element: (
      <RequireAuth>
        <Nav />
        <main>
          <div className="mx-auto max-w-7xl px-2 py-6 sm:px-6 lg:px-8">
            <Domains />
          </div>
        </main>
      </RequireAuth>
    ),
    requrieAuth: true,
  },
  {
    path: "/repo_takedowns",
    name: "Repo Takedowns",
    element: (
      <RequireAuth>
        <Nav />
        <main>
          <div className="mx-auto max-w-7xl px-2 py-6 sm:px-6 lg:px-8">
            <Repos />
          </div>
        </main>
      </RequireAuth>
    ),
    requrieAuth: true,
  },
  {
    path: "/login",
    name: "Login",
    element: (
      <>
        <Nav />
        <main>
          <div className="mx-auto max-w-7xl px-2 py-6 sm:px-6 lg:px-8">
            <Login />
          </div>
        </main>
      </>
    ),
    requrieAuth: false,
    hideIfAuth: true,
  },
  {
    path: "/logout",
    name: "Logout",
    element: (
      <>
        <Nav />
        <main>
          <div className="mx-auto max-w-7xl py-6 px-2 sm:px-6 lg:px-8">
            <Logout />
          </div>
        </main>
      </>
    ),
    requrieAuth: true,
  },
];

const router = createBrowserRouter(routes, {
  basename: "/dash",
});

function Nav() {
  const isAuthed = !!localStorage.getItem("admin_route_token");
  return (
    <Disclosure as="nav" className="bg-gray-800">
      {({ open }) => (
        <>
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="flex h-16 items-center justify-between">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <img
                    className="h-8 w-8"
                    src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=500"
                    alt="BGS Admin Dashboard"
                  />
                </div>
                <div className="hidden md:block">
                  <div className="ml-10 flex items-baseline space-x-4">
                    {routes.map((item) =>
                      (isAuthed && item.hideIfAuth) ||
                      (!isAuthed && item.requrieAuth) ? null : (
                        <NavLink
                          key={item.path}
                          to={item.path || "/"}
                          className={({ isActive }) =>
                            classNames(
                              isActive
                                ? "bg-gray-900 text-white"
                                : "text-gray-300 hover:bg-gray-700 hover:text-white",
                              "rounded-md px-3 py-2 text-sm font-medium"
                            )
                          }
                          aria-current={
                            router.state.location.pathname === item.path
                              ? "page"
                              : undefined
                          }
                        >
                          {item.name}
                        </NavLink>
                      )
                    )}
                  </div>
                </div>
              </div>
              <div className="hidden md:block">
                <div className="ml-4 flex items-center md:ml-6"></div>
              </div>
              <div className="-mr-2 flex md:hidden">
                {/* Mobile menu button */}
                <Disclosure.Button className="inline-flex items-center justify-center rounded-md bg-gray-800 p-2 text-gray-400 hover:bg-gray-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                  <span className="sr-only">Open main menu</span>
                  {open ? (
                    <XMarkIcon className="block h-6 w-6" aria-hidden="true" />
                  ) : (
                    <Bars3Icon className="block h-6 w-6" aria-hidden="true" />
                  )}
                </Disclosure.Button>
              </div>
            </div>
          </div>

          <Disclosure.Panel className="md:hidden">
            <div className="space-y-1 px-2 pb-3 pt-2 sm:px-3">
              {routes.map((item) =>
                (isAuthed && item.hideIfAuth) ||
                (!isAuthed && item.requrieAuth) ? null : (
                  <Disclosure.Button
                    key={item.path}
                    className={classNames(
                      router.state.location.pathname === item.path
                        ? "bg-gray-900 text-white"
                        : "text-gray-300 hover:bg-gray-700 hover:text-white",
                      "block rounded-md px-3 py-2 text-base font-medium"
                    )}
                  >
                    <NavLink
                      key={item.path}
                      to={item.path || "/"}
                      className={({ isActive }) =>
                        classNames(
                          isActive
                            ? "bg-gray-900 text-white"
                            : "text-gray-300 hover:bg-gray-700 hover:text-white",
                          "rounded-md px-3 py-2 text-sm font-medium"
                        )
                      }
                      aria-current={
                        router.state.location.pathname === item.path
                          ? "page"
                          : undefined
                      }
                    >
                      {item.name}
                    </NavLink>
                  </Disclosure.Button>
                )
              )}
            </div>
          </Disclosure.Panel>
        </>
      )}
    </Disclosure>
  );
}

function App() {
  return (
    <>
      <div className="min-h-full">
        <RouterProvider router={router} />
      </div>
    </>
  );
}

export default App;
