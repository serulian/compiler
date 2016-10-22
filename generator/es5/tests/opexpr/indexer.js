$module('indexer', function () {
  var $static = this;
  this.$class('272ade03', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.result = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    $instance.$index = function (someParam) {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $promise.resolve($this.result.$wrapped).then(function ($result0) {
                $result = $t.fastbox($result0 && !someParam.$wrapped, $g.____testlib.basictypes.Boolean);
                $current = 1;
                $continue($resolve, $reject);
                return;
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 1:
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
    $instance.$setindex = function (index, value) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $this.result = value;
        $resolve();
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.indexer.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            sc.$setindex($t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            sc.$index($t.fastbox(false, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
