$module('asyncnullaccess', function () {
  var $static = this;
  this.$class('b0e74b83', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeProperty = $t.property($t.markpromising(function () {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $promise.translate($g.asyncnullaccess.DoSomethingAsync()).then(function ($result0) {
                $result = $t.fastbox($result0.$wrapped == 42, $g.____testlib.basictypes.Boolean);
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
    }));
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProperty|3|9706e8ab": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('69732a26', function () {
    return $t.fastbox(42, $g.____testlib.basictypes.Integer);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            sc = $g.asyncnullaccess.SomeClass.new();
            $t.dynamicaccess(sc, 'SomeProperty', true).then(function ($result0) {
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
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
