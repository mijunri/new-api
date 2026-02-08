/*
Copyright (C) 2025 FoxRouter

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import React, { useContext, useEffect, useState, useRef } from 'react';
import { Button } from '@douyinfe/semi-ui';
import { API, showError, copy, showSuccess } from '../../helpers';
import { useIsMobile } from '../../hooks/common/useIsMobile';
import { StatusContext } from '../../context/Status';
import { useActualTheme } from '../../context/Theme';
import { marked } from 'marked';
import { useTranslation } from 'react-i18next';
import { IconCopy, IconArrowRight } from '@douyinfe/semi-icons';
import { Link } from 'react-router-dom';
import NoticeModal from '../../components/layout/NoticeModal';

// 打字机效果组件
const TypeWriter = ({ texts, speed = 80, deleteSpeed = 40, pauseTime = 2000 }) => {
  const [displayText, setDisplayText] = useState('');
  const [textIndex, setTextIndex] = useState(0);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    const currentText = texts[textIndex];
    
    const timeout = setTimeout(() => {
      if (!isDeleting) {
        if (displayText.length < currentText.length) {
          setDisplayText(currentText.slice(0, displayText.length + 1));
        } else {
          setTimeout(() => setIsDeleting(true), pauseTime);
        }
      } else {
        if (displayText.length > 0) {
          setDisplayText(displayText.slice(0, -1));
        } else {
          setIsDeleting(false);
          setTextIndex((prev) => (prev + 1) % texts.length);
        }
      }
    }, isDeleting ? deleteSpeed : speed);

    return () => clearTimeout(timeout);
  }, [displayText, isDeleting, textIndex, texts, speed, deleteSpeed, pauseTime]);

  return (
    <span className="typewriter-text">
      {displayText}
      <span className="typewriter-cursor">|</span>
    </span>
  );
};

// 代码行动画组件
const CodeLine = ({ delay, children }) => {
  const [visible, setVisible] = useState(false);
  
  useEffect(() => {
    const timer = setTimeout(() => setVisible(true), delay);
    return () => clearTimeout(timer);
  }, [delay]);

  return (
    <div className={`code-line ${visible ? 'visible' : ''}`}>
      {children}
    </div>
  );
};

// 浮动粒子组件
const FloatingParticles = () => {
  return (
    <div className="particles-container">
      {[...Array(20)].map((_, i) => (
        <div
          key={i}
          className="particle"
          style={{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 5}s`,
            animationDuration: `${15 + Math.random() * 10}s`,
          }}
        />
      ))}
    </div>
  );
};

const Home = () => {
  const { t, i18n } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const actualTheme = useActualTheme();
  const [homePageContentLoaded, setHomePageContentLoaded] = useState(false);
  const [homePageContent, setHomePageContent] = useState('');
  const [noticeVisible, setNoticeVisible] = useState(false);
  const isMobile = useIsMobile();
  const serverAddress =
    statusState?.status?.server_address || `${window.location.origin}`;

  const typewriterTexts = [
    'GPT-5.3-Codex',
    'Claude Opus 4.6',
    'Gemini 3 Pro',
  ];

  const displayHomePageContent = async () => {
    setHomePageContent(localStorage.getItem('home_page_content') || '');
    const res = await API.get('/api/home_page_content');
    const { success, message, data } = res.data;
    if (success) {
      let content = data;
      if (!data.startsWith('https://')) {
        content = marked.parse(data);
      }
      setHomePageContent(content);
      localStorage.setItem('home_page_content', content);

      if (data.startsWith('https://')) {
        const iframe = document.querySelector('iframe');
        if (iframe) {
          iframe.onload = () => {
            iframe.contentWindow.postMessage({ themeMode: actualTheme }, '*');
            iframe.contentWindow.postMessage({ lang: i18n.language }, '*');
          };
        }
      }
    } else {
      showError(message);
      setHomePageContent('加载首页内容失败...');
    }
    setHomePageContentLoaded(true);
  };

  const handleCopyBaseURL = async () => {
    const ok = await copy(serverAddress);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };

  useEffect(() => {
    const checkNoticeAndShow = async () => {
      const lastCloseDate = localStorage.getItem('notice_close_date');
      const today = new Date().toDateString();
      if (lastCloseDate !== today) {
        try {
          const res = await API.get('/api/notice');
          const { success, data } = res.data;
          if (success && data && data.trim() !== '') {
            setNoticeVisible(true);
          }
        } catch (error) {
          console.error('获取公告失败:', error);
        }
      }
    };

    checkNoticeAndShow();
  }, []);

  useEffect(() => {
    displayHomePageContent().then();
  }, []);

  return (
    <div className='w-full overflow-x-hidden'>
      <NoticeModal
        visible={noticeVisible}
        onClose={() => setNoticeVisible(false)}
        isMobile={isMobile}
      />
      {homePageContentLoaded && homePageContent === '' ? (
        <div className='claude-home'>
          {/* 背景效果 */}
          <div className='claude-bg'>
            <div className='gradient-orb orb-1' />
            <div className='gradient-orb orb-2' />
            <div className='gradient-orb orb-3' />
            <div className='grid-overlay' />
            <FloatingParticles />
          </div>

          {/* 主要内容区域 */}
          <div className='claude-content'>
            {/* Hero 区域 */}
            <section className='hero-section'>
              <div className='hero-badge'>
                <span className='badge-dot' />
                <span>Claude Code API Proxy</span>
              </div>

              <h1 className='hero-title'>
                <span className='title-line'>为 <span className='gradient-text'>Claude Code</span> 而生</span>
                <span className='title-line'>更快、更稳定、更便宜</span>
              </h1>

              <p className='hero-description'>
                专业的 Claude API 代理服务，完美支持 Claude Code 编程助手
                <br />
                无需翻墙，即刻体验 AI 编程的魅力
              </p>

              {/* 终端样式展示 */}
              <div className='terminal-container'>
                <div className='terminal-header'>
                  <div className='terminal-dots'>
                    <span className='dot red' />
                    <span className='dot yellow' />
                    <span className='dot green' />
                  </div>
                  <span className='terminal-title'>claude-code</span>
                </div>
                <div className='terminal-body'>
                  <CodeLine delay={0}>
                    <span className='code-prompt'>$</span>
                    <span className='code-command'> export ANTHROPIC_BASE_URL=</span>
                    <span className='code-string'>"{serverAddress}"</span>
                  </CodeLine>
                  <CodeLine delay={400}>
                    <span className='code-prompt'>$</span>
                    <span className='code-command'> claude</span>
                  </CodeLine>
                  <CodeLine delay={800}>
                    <span className='code-output'>
                      <span className='code-success'>✓</span> Connected to FoxRouter API
                    </span>
                  </CodeLine>
                  <CodeLine delay={1200}>
                    <span className='code-output'>
                      <span className='code-info'>⚡</span> Model: <TypeWriter texts={typewriterTexts} />
                    </span>
                  </CodeLine>
                  <CodeLine delay={1600}>
                    <span className='code-output'>
                      <span className='code-success'>✓</span> Ready for coding...
                    </span>
                  </CodeLine>
                </div>
              </div>

              {/* 操作按钮 */}
              <div className='hero-actions'>
                <Link to='/console'>
                  <Button
                    theme='solid'
                    size={isMobile ? 'default' : 'large'}
                    className='primary-btn'
                    icon={<IconArrowRight />}
                    iconPosition='right'
                  >
                    开始使用
                  </Button>
                </Link>
                <Button
                  size={isMobile ? 'default' : 'large'}
                  className='secondary-btn'
                  icon={<IconCopy />}
                  onClick={handleCopyBaseURL}
                >
                  复制 API 地址
                </Button>
              </div>
            </section>

            {/* 特性展示 */}
            <section className='features-section'>
              <div className='features-grid'>
                <div className='feature-card'>
                  <div className='feature-icon'>⚡</div>
                  <h3>极速响应</h3>
                  <p>全球 CDN 加速，低延迟高可用</p>
                </div>
                <div className='feature-card'>
                  <div className='feature-icon'>🔒</div>
                  <h3>安全稳定</h3>
                  <p>企业级安全防护，99.9% 可用性</p>
                </div>
                <div className='feature-card'>
                  <div className='feature-icon'>💰</div>
                  <h3>超低价格</h3>
                  <p>按量付费，比官方更优惠</p>
                </div>
                <div className='feature-card'>
                  <div className='feature-icon'>🚀</div>
                  <h3>即插即用</h3>
                  <p>只需修改 Base URL，无需改动代码</p>
                </div>
              </div>
            </section>

          </div>
        </div>
      ) : (
        <div className='overflow-x-hidden w-full'>
          {homePageContent.startsWith('https://') ? (
            <iframe
              src={homePageContent}
              className='w-full h-screen border-none'
            />
          ) : (
            <div
              className='mt-[60px]'
              dangerouslySetInnerHTML={{ __html: homePageContent }}
            />
          )}
        </div>
      )}
    </div>
  );
};

export default Home;
