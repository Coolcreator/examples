from books.models import Book
from rest_framework import viewsets
from user.models import User
from user.serializers import UserSerializer
from books.models import Author, Book
from books.serializers import BookSerializer, AuthorSerializer


class UserModelViewSet(viewsets.ModelViewSet):
    serializer_class = UserSerializer
    queryset = User.objects.all()

class AuthorModelViewSet(viewsets.ModelViewSet):
    serializer_class = AuthorSerializer
    queryset = Author.objects.all()

class BookModelViewSet(viewsets.ModelViewSet):
    serializer_class = BookSerializer
    queryset = Book.objects.all()
    