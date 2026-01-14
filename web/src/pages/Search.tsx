import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search as SearchIcon, File, Folder, Clock, X } from 'lucide-react';
import { Input, Card, Select, Button } from '../components/ui';
import { searchApi } from '../services/api';
import type { File as FileType } from '../types/api';

type SearchIn = 'name' | 'path' | 'all';

export default function Search() {
  const navigate = useNavigate();
  const [query, setQuery] = useState('');
  const [searchIn, setSearchIn] = useState<SearchIn>('all');
  const [results, setResults] = useState<FileType[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (query.trim()) {
        performSearch();
      } else {
        setResults([]);
        setHasSearched(false);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [query, searchIn]);

  const performSearch = async () => {
    if (!query.trim()) return;

    setIsSearching(true);
    setHasSearched(true);

    try {
      const response = await searchApi.searchFiles({
        q: query,
        search_in: searchIn,
        page: 1,
        page_size: 50,
      });
      setResults(response.files);
    } catch (error) {
      console.error('Search failed:', error);
      setResults([]);
    } finally {
      setIsSearching(false);
    }
  };

  const handleSearch = () => {
    performSearch();
  };

  const handleClear = () => {
    setQuery('');
    setResults([]);
    setHasSearched(false);
  };

  const handleFileClick = (file: FileType) => {
    if (file.type === 'directory') {
      navigate('/');
    } else {
      // Open file preview or download
      navigate('/');
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const getFileIcon = (file: FileType) => {
    if (file.type === 'directory') {
      return <Folder className="w-5 h-5 text-blue-500" />;
    }

    const ext = file.name.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return <File className="w-5 h-5 text-red-500" />;
      case 'doc':
      case 'docx':
        return <File className="w-5 h-5 text-blue-600" />;
      case 'xls':
      case 'xlsx':
        return <File className="w-5 h-5 text-green-600" />;
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
        return <File className="w-5 h-5 text-purple-500" />;
      default:
        return <File className="w-5 h-5 text-gray-500" />;
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Search Files</h1>
        <p className="text-sm text-gray-600">
          {hasSearched
            ? `Found ${results.length} result${results.length === 1 ? '' : 's'}`
            : 'Search across all your files'}
        </p>
      </div>

      {/* Search Bar */}
      <Card padding="medium">
        <div className="flex gap-3">
          <div className="flex-1 relative">
            <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <Input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              placeholder="Search for files and folders..."
              className="pl-10 pr-10"
            />
            {query && (
              <button
                onClick={handleClear}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>

          <Select
            value={searchIn}
            onChange={(value) => setSearchIn(value as SearchIn)}
            options={[
              { value: 'all', label: 'All' },
              { value: 'name', label: 'Name' },
              { value: 'path', label: 'Path' },
            ]}
            className="w-32"
          />

          <Button variant="primary" onClick={handleSearch} disabled={isSearching || !query.trim()}>
            Search
          </Button>
        </div>
      </Card>

      {/* Results */}
      {isSearching ? (
        <div className="text-center py-12 text-gray-600">Searching...</div>
      ) : hasSearched && results.length === 0 ? (
        <div className="text-center py-12">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gray-100 rounded-full mb-4">
            <SearchIcon className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No results found</h3>
          <p className="text-gray-600">
            Try adjusting your search terms or filters
          </p>
        </div>
      ) : results.length > 0 ? (
        <Card padding="medium">
          <div className="space-y-2">
            {results.map((file) => (
              <div
                key={file.id}
                onClick={() => handleFileClick(file)}
                className="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
              >
                <div className="flex-shrink-0">{getFileIcon(file)}</div>

                <div className="flex-1 min-w-0">
                  <p className="font-medium text-gray-900 truncate">{file.name}</p>
                  <div className="flex items-center gap-3 mt-1 text-xs text-gray-600">
                    <span>{file.type === 'directory' ? 'Folder' : formatFileSize(file.size)}</span>
                    <span className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      {formatDate(file.updated_at)}
                    </span>
                    <span className="truncate text-gray-500">{file.path}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </Card>
      ) : null}

      {/* Recent Searches (Optional) */}
      {!hasSearched && !query && (
        <Card padding="medium">
          <h3 className="font-medium text-gray-900 mb-3">Search Tips</h3>
          <ul className="space-y-2 text-sm text-gray-600">
            <li>• Use specific keywords to find files faster</li>
            <li>• Search looks in file names by default</li>
            <li>• Select "Path" to search in the full file path</li>
            <li>• Results are updated automatically as you type</li>
          </ul>
        </Card>
      )}
    </div>
  );
}
