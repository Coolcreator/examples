import json
from django.shortcuts import get_object_or_404, render
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.core.exceptions import ObjectDoesNotExist
from rest_framework.decorators import api_view
from rest_framework import status
from .models import Author, Book
from .serializers import AuthorSerializer, BookSerializer


def index(request):
    return render(request, 'index.html')


@api_view(['GET'])
def get_books(request):
    books = Book.objects.all()
    result = BookSerializer(books, many=True)
    return JsonResponse({'books': result.data})


@api_view(['GET'])
def get_book(request, book_id):
    book = get_object_or_404(Book, id=book_id)
    result = BookSerializer(book)
    return JsonResponse({'books': result.data})


@api_view(['POST'])
@csrf_exempt
def add_book(request):
    data = json.loads(request.body)
    serializer = BookSerializer(data=data)
    if serializer.is_valid():
        serializer.save()
        return JsonResponse({'status': 'success'}, status=201)
    return JsonResponse({'status': 'not valid data'}, status=400)


@api_view(['PUT'])
@csrf_exempt
def update_book(request, book_id):
    book_item = get_object_or_404(Book, id=book_id)
    data = json.loads(request.body)
    serializer = BookSerializer(book_item, data=data)
    if serializer.is_valid():
        serializer.save()
        return JsonResponse({'status': 'success'}, status=201)
    return JsonResponse({'status': 'not valid data'}, status=400)


@api_view(['DELETE'])
@csrf_exempt
def delete_book(request, book_id):
    book = get_object_or_404(Book, id=book_id)
    book.delete()
    return JsonResponse({'status': 'success'})
