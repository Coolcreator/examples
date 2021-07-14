from django.db import models
from django.db import models


class Author(models.Model):
    name = models.CharField(max_length=200, verbose_name='Имя автора')
    date_of_birth = models.DateField(null=True, verbose_name='Дата рождения')
    residence = models.CharField(null=True, max_length=200, verbose_name='Место рождения')
    email = models.EmailField(verbose_name='Контакты')

    def __str__(self):
        return self.name


class Book(models.Model):
    title = models.CharField(max_length=200, verbose_name='Имя книги')
    description = models.TextField(verbose_name='Описание книги')
    author = models.ForeignKey(Author, on_delete=models.CASCADE, verbose_name='ID автора')
    rating = models.FloatField(default=0.0, verbose_name='Рейтинг книги')

    def __str__(self):
        return self.title
