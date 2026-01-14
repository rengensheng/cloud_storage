import { ThemeProvider } from "./contexts";
import { AuthProvider } from "./contexts/AuthContext";
import { Router } from "./router";
import "./index.css";

function App() {
  return (
    <ThemeProvider defaultTheme="modern-blue">
      <AuthProvider>
        <Router />
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
