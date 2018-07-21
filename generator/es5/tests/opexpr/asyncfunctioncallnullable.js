$module('asyncfunctioncallnullable', function () {
  var $static = this;
  this.$class('543abf8a', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeMethod = $t.markpromising(function () {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        localasyncloop: while (true) {
          switch ($current) {
            case 0:
              $promise.translate($g.asyncfunctioncallnullable.DoSomethingAsync()).then(function ($result0) {
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
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeMethod|2|cf412abd<aa28dc2d>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('46d4c151', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            sc = $g.asyncfunctioncallnullable.SomeClass.new();
            $t.nullableinvoke(sc, 'SomeMethod', true, []).then(function ($result1) {
              return $promise.resolve($result1).then(function ($result0) {
                $result = $t.asyncnullcompare($result0, $t.fastbox(false, $g.________testlib.basictypes.Boolean));
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
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
