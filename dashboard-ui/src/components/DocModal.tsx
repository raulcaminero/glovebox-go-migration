import React, { useState, useEffect } from 'react';
import { X } from 'lucide-react';

interface DocModalProps {
  docId: string | null;
  title: string;
  onClose: () => void;
}

export const DocModal: React.FC<DocModalProps> = ({ docId, title, onClose }) => {
  const [content, setContent] = useState<string>('Loading document...');

  useEffect(() => {
    if (!docId) return;
    fetch(`/api/docs/${docId}`)
      .then((res) => res.text())
      .then((text) => setContent(text))
      .catch((err) => setContent(`Failed to load doc: ${err.message}`));
  }, [docId]);

  if (!docId) return null;

  return (
    <div className="modal-overlay active">
      <div className="modal-card max-w-2xl">
        <div className="modal-header">
          <h3 className="text-cyan text-sm">{title}</h3>
          <button className="btn btn-secondary text-xs" onClick={onClose}>
            <X size={14} /> Close
          </button>
        </div>
        <div className="modal-body">
          <pre className="doc-pre">{content}</pre>
        </div>
      </div>
    </div>
  );
};
