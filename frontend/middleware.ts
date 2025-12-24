import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Rotas públicas que não requerem autenticação
const publicRoutes = ["/login", "/register"];

// Rotas protegidas que requerem autenticação
const protectedRoutes = ["/accounts", "/transactions", "/reports"];

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get("auth_token")?.value;

  // Verificar se é uma rota pública
  const isPublicRoute = publicRoutes.some((route) =>
    pathname.startsWith(route)
  );

  // Verificar se é uma rota protegida
  const isProtectedRoute = protectedRoutes.some((route) =>
    pathname.startsWith(route)
  );

  // Se não tem token e está tentando acessar rota protegida
  if (!token && isProtectedRoute) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("redirect", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // Se tem token e está tentando acessar rota pública (login/register)
  // NÃO redirecionar automaticamente - deixar o componente ProtectedRoute decidir
  // Isso evita loops de redirecionamento
  // if (token && isPublicRoute) {
  //   return NextResponse.redirect(new URL("/", request.url));
  // }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - swagger (swagger docs)
     * - public files (images, etc.)
     */
    "/((?!api|_next/static|_next/image|favicon.ico|swagger|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};

