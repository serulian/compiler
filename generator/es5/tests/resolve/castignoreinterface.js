$module('castignoreinterface', function () {
  var $static = this;
  this.$class('7fd03495', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeFunction = $t.markpromising(function (value) {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $promise.translate($g.castignoreinterface.DoSomethingAsync()).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              }).catch(function (err) {
                $reject(err);
                return;
              });
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
        "SomeFunction|2|fd8bc7c9<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('750abcaf', 'SomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeFunction|2|fd8bc7c9<void>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('7c17decb', function () {
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var a;
    var somevalue;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            somevalue = $t.fastbox('hello', $g.________testlib.basictypes.String);
            try {
              var $expr = $t.cast(somevalue, $g.castignoreinterface.SomeInterface, false);
              a = $expr;
            } catch ($rejected) {
              a = null;
            }
            $current = 1;
            $continue($resolve, $reject);
            return;

          case 1:
            $t.nullableinvoke(a, 'SomeFunction', true, [$t.fastbox(true, $g.________testlib.basictypes.Boolean)]).then(function ($result0) {
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
            $resolve($t.fastbox(a == null, $g.________testlib.basictypes.Boolean));
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
