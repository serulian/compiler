$module('indexer', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.nominalwrap(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.result = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.$index = function (someParam) {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $promise.resolve($t.nominalunwrap($this.result)).then(function ($result0) {
                $result = $t.nominalwrap($result0 && !$t.nominalunwrap(someParam), $g.____testlib.basictypes.Boolean);
                $state.current = 1;
                $callback($state);
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
    $instance.$setindex = function (index, value) {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $this.result = value;
              $state.current = -1;
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
  });

  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.indexer.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            sc.$setindex($t.nominalwrap(1, $g.____testlib.basictypes.Integer), $t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            sc.$index($t.nominalwrap(false, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $state.resolve($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
