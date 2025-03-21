package controllers

import (
	"context"
	"demariot-backend/database"
	"demariot-backend/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(c *fiber.Ctx) error {
	collection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener los usuarios"})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al procesar los datos de los usuarios"})
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los resultados de la base de datos"})

	}

	return c.JSON(users)
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Id invalido"})
	}

	collection := database.GetCollection("users")

	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	return c.JSON(user)
}

func UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var request struct {
		Role string `json:"role"`
	}

	// Parsear el cuerpo de la solicitud
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear la solicitud"})
	}

	// Validar si el rol está vacío
	if request.Role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "El rol no puede estar vacío"})
	}

	// Obtener la colección de usuarios
	collection := database.GetCollection("users")

	// Contexto de la base de datos con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Actualizar el rol en la base de datos
	update := bson.M{"$set": bson.M{"role": request.Role, "updated_at": time.Now()}}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo actualizar el rol del usuario"})
	}

	// Respuesta exitosa
	return c.JSON(fiber.Map{"message": "Rol actualizado correctamente"})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, _ := primitive.ObjectIDFromHex(id)

	collection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.DeleteOne(ctx, bson.M{"_id": objectId})

	return c.JSON(fiber.Map{"message": "Usuario eliminado correctamente"})
}

func UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Evitar modificar el _id
	delete(updateData, "_id")

	// Agregar timestamp de actualización
	updateData["updated_at"] = time.Now()

	update := bson.M{"$set": updateData}
	filter := bson.M{"_id": objId}

	collection := database.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ejecutar la actualización
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar usuario"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	return c.JSON(fiber.Map{"message": "Usuario actualizado correctamente"})
}

func UploadProfilePicture(c *fiber.Ctx) error {
	// Obtener el archivo cargado con FormFile
	file, err := c.FormFile("profile_picture")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No se pudo procesar el archivo"})
	}

	// Crear un nombre único para la imagen
	uploadPath := "./uploads/" + utils.UUIDv4() + ".jpg"

	// Guardar el archivo en la carpeta "uploads"
	err = c.SaveFile(file, uploadPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar la imagen"})
	}

	// Obtener el ID del usuario desde la URL y convertirlo a ObjectID
	userID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de usuario no válido"})
	}

	// Actualizar la base de datos con el nuevo nombre de archivo
	collection := database.GetCollection("users")
	_, err = collection.UpdateOne(c.Context(), bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{"profile_picture": uploadPath},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar la imagen de perfil"})
	}

	// Retornar un mensaje de éxito y la ruta del archivo cargado
	return c.JSON(fiber.Map{
		"message":         "Imagen de perfil cargada con éxito",
		"profile_picture": uploadPath,
	})
}
