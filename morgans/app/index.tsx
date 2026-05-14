import { useEffect, useState } from 'react';
import { View, Text } from 'react-native';
import { Article } from '../types/article';

export default function Index() {
  const [articles, setArticles] = useState<Article[]>([])

  useEffect(() => {
    const fetchArticles = async () => {
      try {
        const response = await fetch('http://egghead:8000/articles');
        const data: Article[] = await response.json();
        // Response handling will be implemented later
        console.log(`Fetched ${data.length} articles`);

        setArticles(data)
      } catch (error) {
        console.error('Error fetching articles:', error);
      }
    };

    fetchArticles();
  }, []);

  return (
    <View className="flex-1 bg-white items-center justify-center">
      <Text className="text-2xl font-bold">Unread</Text>
      {
        articles.map(article => {
          if (!article.isRead) {
            return <Text>{article.title}</Text> 
          }
        })
      }

      <Text className="text-2xl font-bold">Read</Text>
      {
        articles.map(article => {
          if (article.isRead) {
            return <Text>{article.title}</Text> 
          }
        })
      }
    </View>
  );
}
