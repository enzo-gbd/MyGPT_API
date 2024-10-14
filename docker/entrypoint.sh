#!/bin/sh

echo "En attente de la base de données..."
while ! nc -z db 5432; do
  sleep 1
done
echo "Base de données prête."

echo "Exécution des migrations..."
./migration

echo "Démarrage de l'application..."
exec ./GBA
