(function () {
  var seg = window.location.pathname.replace(/^\//, '').split('/')[0];
  var isChapterPage = /^\d{2}-/.test(seg);
  var currentSlug = isChapterPage ? seg : null;
  var previousCompleted = null;
  var bannerShown = false;

  function poll() {
    fetch('/api/progress')
      .then(function (r) { return r.json(); })
      .then(function (data) {
        updateSidebar(data);

        if (!isChapterPage || bannerShown) return;

        var completed = new Set(data.completed || []);

        if (previousCompleted !== null) {
          var currentChapter = null;
          for (var i = 0; i < data.chapters.length; i++) {
            if (data.chapters[i].slug === currentSlug) {
              currentChapter = data.chapters[i];
              break;
            }
          }
          if (!currentChapter) return;

          var wasCompleted = previousCompleted.has(currentChapter.id);
          var isNowCompleted = completed.has(currentChapter.id);

          if (!wasCompleted && isNowCompleted) {
            var idx = data.chapters.indexOf(currentChapter);
            var nextChapter = (idx >= 0 && idx + 1 < data.chapters.length)
              ? data.chapters[idx + 1]
              : null;
            showBanner(currentChapter, nextChapter);
            bannerShown = true;
            return;
          }
        }

        previousCompleted = completed;
      })
      .catch(function () {}); // silently ignore errors
  }

  function updateSidebar(data) {
    var completed = new Set(data.completed || []);
    var unlocked = new Set(data.unlocked || []);

    var links = document.querySelectorAll('a.gdoc-nav__entry');
    for (var i = 0; i < links.length; i++) {
      var link = links[i];
      var href = link.getAttribute('href') || '';
      var match = href.match(/\/(\d{2}-[^/]+)\//);
      if (!match) continue;
      var slug = match[1];

      // find chapter by slug
      var chapter = null;
      for (var j = 0; j < data.chapters.length; j++) {
        if (data.chapters[j].slug === slug) {
          chapter = data.chapters[j];
          break;
        }
      }
      if (!chapter) continue;

      var icon;
      if (completed.has(chapter.id)) {
        icon = '\u2705'; // ✅
      } else if (unlocked.has(chapter.id)) {
        icon = '\uD83D\uDD13'; // 🔓
      } else {
        icon = '\uD83D\uDD12'; // 🔒
      }

      var span = link.querySelector('.ws-status-icon');
      if (!span) {
        span = document.createElement('span');
        span.className = 'ws-status-icon';
        link.appendChild(span);
      }
      span.textContent = icon;
    }
  }

  function showBanner(current, next) {
    var banner = document.createElement('div');
    banner.style.cssText = [
      'position:fixed', 'bottom:0', 'left:0', 'right:0', 'z-index:9999',
      'background:#313244', 'color:#cdd6f4',
      'padding:1rem 1.5rem',
      'display:flex', 'align-items:center', 'justify-content:space-between',
      'box-shadow:0 -4px 16px rgba(0,0,0,0.4)',
      'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif',
      'font-size:0.95rem',
      'gap:1rem',
    ].join(';');

    var msg = document.createElement('span');
    if (next) {
      msg.innerHTML = '<strong style="color:#a6e3a1">\u2713 Chapter ' + current.id + ' complete!</strong>'
        + '&nbsp; Chapter ' + next.id + ' (<em>' + next.name + '</em>) is now unlocked.';
    } else {
      msg.innerHTML = '<strong style="color:#a6e3a1">\u2713 Chapter ' + current.id + ' complete!</strong>'
        + '&nbsp; You have completed the entire workshop!';
    }

    var right = document.createElement('div');
    right.style.cssText = 'display:flex;align-items:center;gap:1rem;flex-shrink:0';

    if (next) {
      var link = document.createElement('a');
      link.href = '/' + next.slug + '/';
      link.textContent = 'Continue to Chapter ' + next.id + ' \u2192';
      link.style.cssText = 'color:#89b4fa;text-decoration:none;font-weight:600;white-space:nowrap';
      right.appendChild(link);
    }

    var close = document.createElement('button');
    close.textContent = '\u00d7';
    close.setAttribute('aria-label', 'Close');
    close.style.cssText = 'background:none;border:none;color:#a6adc8;font-size:1.4rem;cursor:pointer;padding:0 0.25rem;line-height:1';
    close.addEventListener('click', function () { banner.remove(); });
    right.appendChild(close);

    banner.appendChild(msg);
    banner.appendChild(right);
    document.body.appendChild(banner);
  }

  function fetchAndRenderChecks() {
    var container = document.getElementById('ws-check-results');
    if (!container) return;
    var chapterID = container.getAttribute('data-chapter-id');
    if (!chapterID) return;

    fetch('/api/checks')
      .then(function (r) { return r.json(); })
      .then(function (store) {
        var results = store[chapterID];
        if (!results || results.length === 0) return;

        var html = '<h2>Last Check Results</h2><ul class="ws-checklist">';
        for (var i = 0; i < results.length; i++) {
          var r = results[i];
          var cls = r.passed ? 'passed' : 'failed';
          var icon = r.passed ? '\u2713' : '\u2717';
          html += '<li class="ws-check-item ' + cls + '">'
            + '<div class="ws-check-row">'
            + '<span class="ws-check-icon">' + icon + '</span>'
            + '<span class="ws-check-desc">' + r.description + '</span>'
            + '</div>';
          if (!r.passed && r.error) {
            html += '<span class="ws-check-error">' + r.error + '</span>';
          }
          html += '</li>';
        }
        html += '</ul>';
        container.innerHTML = html;
        container.style.display = 'block';
      })
      .catch(function () {});
  }

  function pollAll() {
    poll();
    fetchAndRenderChecks();
  }

  // Initial fetch after a short delay.
  setTimeout(pollAll, 500);

  // SSE for instant updates when check is run.
  var es = new EventSource('/api/events');
  es.onmessage = function (e) {
    if (e.data === 'update' || e.data === 'connected') {
      pollAll();
    }
  };
  es.onerror = function () {
    es.close();
    // Fall back to polling every 5s if SSE fails.
    setInterval(pollAll, 5000);
  };

  // Keepalive polling every 30s as safety net.
  setInterval(pollAll, 30000);
})();
