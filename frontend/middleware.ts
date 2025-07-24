import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
    const sessionToken = request.cookies.get('session')?.value

    if (!sessionToken) {
        const loginUrl = new URL('/auth/login', request.url)
        return NextResponse.redirect(loginUrl)
    }

    return NextResponse.next()
}

export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api (API routes)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         * - /auth (login/register pages)
         * - / (homepage)
         */
        '/((?!api|_next/static|_next/image|favicon.ico|auth|$).*)',
    ],
}