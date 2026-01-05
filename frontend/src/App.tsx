import React, { useState, useEffect } from 'react';

const API_URL = 'http://localhost:8080';

export default function App() {
  const [longUrl, setLongUrl] = useState('');
  const [urls, setUrls] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  useEffect(() => {
    fetchUrls();
  }, []);

  const fetchUrls = async () => {
    try {
      const response = await fetch(`${API_URL}/api/urls`);
      const data = await response.json();
      setUrls(data || []);
    } catch (err) {
      console.error('Erro ao buscar URLs:', err);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccessMessage('');
    setLoading(true);

    try {
      const response = await fetch(`${API_URL}/api/shorten`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ long_url: longUrl }),
      });

      const data = await response.json();
      setSuccessMessage(`URL encurtada: ${data.short_url}`);
      setLongUrl('');
      fetchUrls();
    } catch (err) {
      setError('Erro ao encurtar URL. Verifique se o backend está rodando.');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (shortCode) => {
    if (!window.confirm('Deseja deletar esta URL?')) return;

    try {
      await fetch(`${API_URL}/api/urls/${shortCode}`, { method: 'DELETE' });
      fetchUrls();
    } catch (err) {
      alert('Erro ao deletar URL');
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    alert('URL copiada!');
  };

  return (
    <div className="min-h-screen w-full bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold text-gray-800 mb-3">
             Encurtador de URLs
          </h1>
          <p className="text-gray-600 text-lg">
            Transforme URLs longas em links curtos e rastreáveis
          </p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8 mb-8">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Cole sua URL longa aqui:
              </label>
              <input
                type="url"
                value={longUrl}
                onChange={(e) => setLongUrl(e.target.value)}
                placeholder="https://exemplo.com/minha-url-muito-longa"
                required
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition"
              />
            </div>

            <button
              onClick={handleSubmit}
              disabled={loading || !longUrl}
              className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-6 rounded-lg transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? ' Encurtando...' : ' Encurtar URL'}
            </button>
          </div>

          {error && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
               {error}
            </div>
          )}

          {successMessage && (
            <div className="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg">
              <p className="text-green-800 font-medium"> {successMessage}</p>
            </div>
          )}
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">
             URLs Criadas ({urls.length})
          </h2>

          {urls.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <p className="text-lg">Nenhuma URL criada ainda</p>
              <p className="text-sm mt-2">Comece encurtando sua primeira URL acima! </p>
            </div>
          ) : (
            <div className="space-y-4">
              {urls.map((url) => (
                <div
                  key={url.id}
                  className="border border-gray-200 rounded-lg p-5 hover:shadow-md transition"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1 min-w-0">
                      <div className="mb-3">
                        <p className="text-sm text-gray-600 mb-1">URL Curta:</p>
                        <div className="flex items-center gap-2">
                          <a
                            href={`${API_URL}/${url.short_code}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-indigo-600 font-semibold hover:text-indigo-800 break-all"
                          >
                            {API_URL}/{url.short_code}
                          </a>
                          <button
                            onClick={() => copyToClipboard(`${API_URL}/${url.short_code}`)}
                            className="flex-shrink-0 text-gray-500 hover:text-indigo-600 transition"
                            title="Copiar"
                          >
                            Copiar
                          </button>
                        </div>
                      </div>

                      <div className="mb-3">
                        <p className="text-sm text-gray-600 mb-1">URL Original:</p>
                        <p className="text-gray-800 break-all text-sm">{url.long_url}</p>
                      </div>

                      <div className="flex items-center gap-6 text-sm text-gray-600">
                        <span className="flex items-center gap-1">
                           <strong>{url.clicks}</strong> cliques
                        </span>
                        <span>
                           {new Date(url.created_at).toLocaleDateString('pt-BR')}
                        </span>
                      </div>
                    </div>

                    <button
                      onClick={() => handleDelete(url.short_code)}
                      className="flex-shrink-0 text-red-500 hover:text-red-700 font-medium px-3 py-2 rounded hover:bg-red-50 transition"
                      title="Deletar"
                    >
                      X
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="text-center mt-8 text-gray-600">
          <p className="text-sm">Feito por Rafael Ramalho Rosa </p>
        </div>
      </div>
    </div>
  );
}