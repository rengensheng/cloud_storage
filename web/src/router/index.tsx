import { Suspense, lazy } from "react";
import { Routes, Route } from "react-router-dom";
import LoadingSpinner from "../components/LoadingSpinner";
import { PublicRoute, ProtectedRoute } from "../components/ProtectedRoute";
import MainLayout from "../layouts/MainLayout";

const Login = lazy(() => import("../pages/Login"));
const Register = lazy(() => import("../pages/Register"));
const Files = lazy(() => import("../pages/Files"));
const Upload = lazy(() => import("../pages/Upload"));
const Search = lazy(() => import("../pages/Search"));
const Recycle = lazy(() => import("../pages/Recycle"));
const Settings = lazy(() => import("../pages/Settings"));
const NotFound = lazy(() => import("../pages/NotFound"));

function LazyComponent({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={<LoadingSpinner message="Loading page..." size="md" />}>
      {children}
    </Suspense>
  );
}

export function Router() {
  return (
    <Routes>
      <Route
        path="/login"
        element={
          <PublicRoute>
            <LazyComponent>
              <Login />
            </LazyComponent>
          </PublicRoute>
        }
      />
      <Route
        path="/register"
        element={
          <PublicRoute>
            <LazyComponent>
              <Register />
            </LazyComponent>
          </PublicRoute>
        }
      />

      <Route
        path="/"
        element={
          <ProtectedRoute>
            <MainLayout />
          </ProtectedRoute>
        }
      >
        <Route
          index
          element={
            <LazyComponent>
              <Files />
            </LazyComponent>
          }
        />
        <Route
          path="upload"
          element={
            <LazyComponent>
              <Upload />
            </LazyComponent>
          }
        />
        <Route
          path="search"
          element={
            <LazyComponent>
              <Search />
            </LazyComponent>
          }
        />
        <Route
          path="recycle"
          element={
            <LazyComponent>
              <Recycle />
            </LazyComponent>
          }
        />
        <Route
          path="settings"
          element={
            <LazyComponent>
              <Settings />
            </LazyComponent>
          }
        />
      </Route>
    </Routes>
  );
}
