import { useEffect, useMemo, useState } from 'react';
import {
  View,
  Text,
  SectionList,
  ActivityIndicator,
  Pressable,
} from 'react-native';
import { Linking } from 'react-native';
import { Article } from '../../types/articles';

const API_URL = 'http://egghead:8000';

export default function IndexScreen() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(`${API_URL}/articles`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch articles');
        return res.json();
      })
      .then((data: Article[] | null) => {
        setArticles(data ?? []);
      })
      .catch((err) => {
        setError(err.message);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  const sections = useMemo(() => {
    const unread = articles.filter((a) => !a.isRead);
    const read = articles.filter((a) => a.isRead);
    const result = [];
    if (unread.length > 0) result.push({ title: 'Unread', data: unread });
    if (read.length > 0) result.push({ title: 'Read', data: read });
    return result;
  }, [articles]);

  if (loading) {
    return (
      <View className="flex-1 items-center justify-center">
        <ActivityIndicator size="large" />
      </View>
    );
  }

  if (error) {
    return (
      <View className="flex-1 items-center justify-center p-4">
        <Text className="text-red-500">{error}</Text>
      </View>
    );
  }

  if (articles.length === 0) {
    return (
      <View className="flex-1 items-center justify-center">
        <Text className="text-gray-500">No articles today</Text>
      </View>
    );
  }

  const handleArticlePress = (article: Article) => {
    Linking.openURL(article.link);

    if (!article.isRead) {
      fetch(`${API_URL}/articles`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ link: article.link, isRead: true }),
      }).then(() => {
        setArticles((prev) =>
          prev.map((a) =>
            a.link === article.link ? { ...a, isRead: true } : a
          )
        );
      });
    }
  };

  return (
    <SectionList
      className="flex-1"
      contentContainerClassName="p-4"
      sections={sections}
      keyExtractor={(item) => item.link}
      renderSectionHeader={({ section }) => (
        <Text className="text-lg font-bold pt-4 pb-2">{section.title}</Text>
      )}
      renderItem={({ item }) => (
        <Pressable
          onPress={() => handleArticlePress(item)}
          className="py-3 border-b border-gray-200"
        >
          <Text className="text-base text-blue-600">{item.title}</Text>
        </Pressable>
      )}
    />
  );
}
