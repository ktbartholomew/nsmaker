(function() {
  var getCSRFToken = function() {
    var match = document.cookie.match(/_csrf=(.+?);?$/);
    if (match === null) {
      return '';
    }

    return match[1];
  };

  var StatusMessage = function() {
    var element = document.createElement('div');
    var state = {
      message: '',
      className: ''
    };

    this.Render = function() {
      element.className = 'form-group ' + state.className;
      element.textContent = state.message;
    }.bind(this);

    this.GetElement = function() {
      return element;
    };

    this.SetMessage = function(message) {
      state.message = message;
      requestAnimationFrame(this.Render);
    };

    this.SetClassName = function(className) {
      state.className = className;
      requestAnimationFrame(this.Render);
    };

    this.Render();
  };

  var status = new StatusMessage();
  var form = document.getElementById('form-create-namespace');
  form.appendChild(status.GetElement());

  form['namespace'].focus();
  form.addEventListener('submit', function(e) {
    e.preventDefault();

    status.SetMessage('');
    status.SetClassName('');

    var postData = {
      namespace: form['namespace'].value,
      username: form['username'].value
    };

    (function(d) {
      var failed = false;
      var statusCode = 0;

      fetch('/create', {
        method: 'POST',
        cache: 'no-cache',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': getCSRFToken()
        },
        body: JSON.stringify(d)
      })
        .then(function(res) {
          if (!res.ok) {
            failed = true;
            statusCode = res.status;
            return res.json();
          }

          return res.text();
        })
        .then(function(body) {
          if (failed) {
            throw new Error(body.message);
          }

          status.SetMessage(`namespace "${form['namespace'].value}" created`);
          status.SetClassName('text-success');

          form['namespace'].value = '';
          form['username'].value = '';
        })
        .catch(function(err) {
          status.SetMessage(err);
          status.SetClassName('text-error');
        });
    })(postData);
  });
})();
