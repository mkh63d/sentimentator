# Multi-stage Dockerfile for Sentimentator.Backend.Api
# Uses .NET 10 SDK to build and ASP.NET runtime to run the app.

FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

# Copy project files and restore to leverage cache
COPY ["Sentimentator.Backend.Api/Sentimentator.Backend.Api.csproj", "Sentimentator.Backend.Api/"]
COPY ["Sentimentator.Backend/Sentimentator.Backend.csproj", "Sentimentator.Backend/"]
COPY ["Sentimentator.Backend.Service/Sentimentator.Backend.Service.csproj", "Sentimentator.Backend.Service/"]

RUN dotnet restore "Sentimentator.Backend.Api/Sentimentator.Backend.Api.csproj"

# Copy entire repo and publish
COPY . .
WORKDIR /src/Sentimentator.Backend.Api
RUN dotnet publish -c Release -o /app/publish --no-restore /p:UseAppHost=false

# Runtime image
FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS runtime
WORKDIR /app
COPY --from=build /app/publish .

EXPOSE 80

ENTRYPOINT ["dotnet", "Sentimentator.Backend.Api.dll"]
