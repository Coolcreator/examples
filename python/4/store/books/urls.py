from django.urls import path
from . import views
from .viewssets import *
from rest_framework.routers import DefaultRouter

urlpatterns = [
    path('', views.index),
    path('books', views.get_books),
    path('addbook', views.add_book),
    path('book/<int:book_id>', views.get_book),
    path('updatebook/<int:book_id>', views.update_book),
    path('deletebook/<int:book_id>', views.delete_book),
]

router = DefaultRouter()
router.register(r'users', UserModelViewSet, basename='user')
router.register(r'authors/views', AuthorModelViewSet, basename='author')
router.register(r'books/views', BookModelViewSet, basename='book')

urlpatterns += router.urls
