# ASP.NET Core Identity

This project demonstrates the use of ASP.NET Core Identity.

## Quick Start

### Configure

Update the `ConnectionStrings:DefaultConnection` configuration at `appsettings.Development.json`.

### Register

Register a new user:

```
POST /register
{
  "email": "string",
  "password": "string"
}
```

### Login

Login with `useCookies:true` and `useSessionCookies:true`:

```
POST /login
{
  "email": "string",
  "password": "string"
}
```

### Get

Get login data:

```
GET /manage/info
```