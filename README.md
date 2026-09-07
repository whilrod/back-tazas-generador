# 🖥️ Generador de Plantillas - Backend  

## API REST construida en **Go (Golang)**, desplegada en **Render**, que gestiona imágenes, genera PDFs y se conecta a la base de datos **Xata/Postgres**.  

---

### 🔹 Tipo de arquitectura  

👉 Es parte de una arquitectura distribuida basada en microservicios / servicios independientes:  

- **Frontend (Netlify)** → Cliente SPA en React que consume esta API.  
- **Backend (Render)** → Servicio REST en Go.  
- **Database (Supabase)** → Base de datos gestionada como servicio independiente (DBaaS).  
- **Storage (Cloudinary/Xata files)** → Repositorio de imágenes externas.  

📌 Modelo: **Backend desacoplado (BFF)**  
- Responde únicamente a las necesidades del frontend.  
- Exposición de endpoints REST simples: `/images`, `/images/hashtag`, `/images/pdf`.  

**Patrón de diseño aplicado:**  

- **Backend for Frontend (BFF)**: API pensada solo para este frontend.  
- **Repository Pattern**: acceso a la base de datos encapsulado en un módulo (`db`).  
- **Separation of Concerns (SoC)**: Handlers, DB y main separados en capas.  

---

## 🚀 Características  
- 📦 API REST para listar imágenes y buscarlas por hashtags  
- 🖼️ Descarga de imágenes desde URLs externas (Cloudinary / Storage)  
- 📄 Generación de PDFs dinámicos con las imágenes seleccionadas  
- 🔒 Configuración de **CORS** para permitir consumo desde el frontend (Netlify)  
- 🌍 Desplegado en **Render** con autoscaling y HTTPS  

---

## 📂 Estructura del proyecto  
![alt text](image.png)

## ⚙️ Instalación local

### 1. Clonar el repositorio:

git clone https://github.com/tuusuario/back-tazas-generador.git
cd back-tazas-generador


### 2. Instalar dependencias:

go mod tidy


### 3. Configurar variables de entorno en .env:

PORT=8080
ALLOWED_ORIGIN=http://localhost:5173
XATA_DATABASE_URL=postgres://usuario:password@host:5432/dbname

### 4. Ejecutar servidor local:

go run main.go

👉 Servidor en: http://localhost:8080

## 🌍 Despliegue en Render

### Build Command

go build -tags netgo -ldflags '-s -w' -o app


### Start Command

./app


### Variables de entorno en Render

PORT=10000
ALLOWED_ORIGIN=https://plantillastazas.netlify.app
XATA_DATABASE_URL=postgres://usuario:password@host:5432/dbname


## 👉 Endpoint desplegado:
https://back-tazas-generador.onrender.com

📑 Endpoints principales
🔹 Listar imágenes
GET /images?page=1&limit=20

🔹 Buscar por hashtags
GET /images/hashtag?tag=taza&tag=regalo&page=1&limit=10

🔹 Generar PDF
POST /images/pdf
Body:
{
  "uuids": ["uuid1", "uuid2", "uuid3"]
}

![alt text](image-1.png) ![alt text](image-2.png) ![alt text](image-3.png)
![alt text](image-4.png)
## 📄 Respuesta: PDF descargable con 3 imágenes por página.

🛠️ Tecnologías utilizadas

Go (Golang)

gofpdf (para generación de PDFs)

pq (driver PostgreSQL)

Xata / Postgres como DBaaS

Render (hosting backend)

rs/cors (middleware CORS)
